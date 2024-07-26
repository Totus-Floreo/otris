// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This library has been modified by Timur Kulakov for the open-source project Open Streaming Solutions in 2024.

package otris

import (
	"bytes"
	"context"
	"github.com/Totus-Floreo/otris/internal/slog/buffer"
	"io"
	"log/slog"
	"slices"
	"sync"
)

// Handler is a modified version of the original commonHandler from the log/slog.
type Handler struct {
	// commonHandler fields
	json              bool                 // Default for json is false
	pretty            bool                 // Default for pretty is false
	safe              bool                 // The `safe` field is a boolean flag that indicates whether the handler is in a safe set or not.
	sep               string               // Default for sep is " "
	layout            string               // Default for layout is otris.DefaultDateTimeLayout
	color             LevelColorMap        // Color map for different log levels
	opts              *slog.HandlerOptions // Warning! HandlerOptions is WIP in v2. You can use it, but at one's own risk.
	preformattedAttrs []byte
	groupPrefix       string
	groups            []string
	nOpenGroups       int
	buf               *bytes.Buffer
	mu                *sync.Mutex
	w                 io.Writer
}

// NewHandler is manually constructor, please use NewHandlerBuilder.
// `safeSet` from coding/json/tables.go is used to escape the string
// Warning! HandlerOptions is WIP. You can use it, but at one's own risk.
// Only for tests, use builder please! Default setting close to NewPrettyHandler.
func NewHandler(w io.Writer, color LevelColorMap, safe bool, layout string, sep string, opts *slog.HandlerOptions) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	if layout == "" {
		layout = DefaultPrettyDateTimeLayout
	}
	if sep == "" {
		sep = PrettySep
	}
	if color == nil {
		color = DefaultColorMap
	}
	return &Handler{
		json:   false,
		pretty: true,
		safe:   safe,
		color:  color,
		layout: layout,
		sep:    sep,
		w:      w,
		opts:   opts,
		mu:     &sync.Mutex{},
	}
}

// NewPrettyHandler is default otris handler without settings.
// If you need advanced settings please use NewHandlerBuilder
func NewPrettyHandler(w io.Writer, opts *slog.HandlerOptions) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &Handler{
		json:   false,
		pretty: true,
		safe:   false,
		color:  DefaultColorMap,
		layout: DefaultPrettyDateTimeLayout,
		sep:    PrettySep,
		w:      w,
		opts:   opts,
		mu:     &sync.Mutex{},
	}
}

// NewJSONHandler is analog for slog.NewJSONHandler
func NewJSONHandler(w io.Writer, opts *slog.HandlerOptions) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &Handler{
		json:   true,
		pretty: false,
		safe:   true,
		color:  EmptyColorMap,
		sep:    JSONSep,
		w:      w,
		opts:   opts,
		mu:     &sync.Mutex{},
	}
}

// NewStructHandler is analog for slog.NewTextHandler
func NewStructHandler(w io.Writer, opts *slog.HandlerOptions) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &Handler{
		json:   false,
		pretty: false,
		safe:   true,
		sep:    StructSep,
		w:      w,
		opts:   opts,
		mu:     &sync.Mutex{},
	}
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	minLevel := LevelFx
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

func (h *Handler) Handle(_ context.Context, record slog.Record) error {
	// Use an empty separator for reuse later, since it is always inserted during state.append...
	state := h.newHandleState(buffer.New(), true, "")
	defer state.free()
	if h.json {
		state.buf.WriteByte('{')
	}
	// Built-in attributes. They are not in a group.
	stateGroups := state.groups
	state.groups = nil // So ReplaceAttrs sees no groups instead of the pre groups.
	rep := h.opts.ReplaceAttr

	// time
	if !record.Time.IsZero() {
		key := slog.TimeKey
		val := record.Time.Round(0) // strip monotonic to match Attr behavior
		if rep == nil {
			state.appendKey(key)
			state.appendTime(val)
		} else {
			state.appendAttr(slog.Time(key, val)) // <- TODO Refactor state.appendAttr in v2
		}
	}

	// level
	key := slog.LevelKey
	val := record.Level
	state.color = GetColor(h.color, val)
	if rep == nil {
		state.appendKey(key)
		state.appendString(GetLevelName(val))
	} else {
		state.appendAttr(slog.Any(key, val)) // <- TODO Refactor state.appendAttr in v2
	}
	state.resetColor()

	// source
	if h.opts.AddSource {
		state.appendAttr(slog.Any(slog.SourceKey, rSource(record))) // <- TODO Refactor state.appendAttr in v2
	}
	key = slog.MessageKey
	msg := record.Message
	if rep == nil {
		state.appendKey(key)
		state.appendString(msg)
	} else {
		state.appendAttr(slog.String(key, msg)) // <- TODO Refactor state.appendAttr in v2
	}
	state.groups = stateGroups // Restore groups passed to ReplaceAttrs.
	state.appendNonBuiltIns(record)
	state.buf.WriteByte('\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.Write(*state.buf)
	return err
}

// WithAttrs returns a new Handler with additional attributes specified in `attrs` parameter.
// TODO Implement custom groupPrefix in v2
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// We are going to ignore empty groups, so if the entire slice consists of
	// them, there is nothing to do.
	if countEmptyGroups(attrs) == len(attrs) {
		return h
	}
	h2 := h.clone()
	// Pre-format the attributes as an optimization.
	// Use an empty separator for reuse later, since it is always inserted during state.append...
	state := h2.newHandleState((*buffer.Buffer)(&h2.preformattedAttrs), false, "")
	defer state.free()
	state.prefix.WriteString(h.groupPrefix)
	if len(h2.preformattedAttrs) > 0 {
		state.sep = h.attrSep()
	}
	state.openGroups()
	for _, a := range attrs {
		state.appendAttr(a)
	}
	// Remember the new prefix for later keys.
	h2.groupPrefix = state.prefix.String()
	// Remember how many opened groups are in preformattedAttrs,
	// so we don't open them again when we handle a Record.
	h2.nOpenGroups = len(h2.groups)
	return h2
}

// WithGroup returns a new Handler with the given `name` appended to the `groups` slice.
func (h *Handler) WithGroup(name string) slog.Handler {
	h2 := h.clone()
	h2.groups = append(h2.groups, name)
	return h2
}

func (h *Handler) clone() *Handler {
	return &Handler{
		json:              h.json,
		opts:              h.opts,
		preformattedAttrs: slices.Clip(h.preformattedAttrs),
		groupPrefix:       h.groupPrefix,
		groups:            slices.Clip(h.groups),
		nOpenGroups:       h.nOpenGroups,
		w:                 h.w,
		mu:                h.mu,
	}
}

// attrSep returns the separator between attributes.
func (h *Handler) attrSep() string {
	// use a boolean json to avoid unnecessary errors
	if h.json {
		return JSONSep
	}
	return h.sep
}
