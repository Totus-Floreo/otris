package otris

import (
	"io"
	"log/slog"
	"os"
	"sync"
)

// HandlerBuilder is a type that helps in building a Handler by setting various options.
//
// Usage:
//
//	handler := NewHandlerBuilder().WithColor(color).WithSafeSet(safe).WithTimeLayout(layout).Build()
type HandlerBuilder struct {
	h *Handler
}

// NewHandlerBuilder creates a new instance of HandlerBuilder. It initializes the fields of HandlerBuilder
// with default values. The returned HandlerBuilder can be used to configure and build a Handler.
//
// Returns a pointer to the created HandlerBuilder.
func NewHandlerBuilder() *HandlerBuilder {
	return &HandlerBuilder{
		h: &Handler{
			json:   false,
			pretty: false,
			safe:   true,
			color:  EmptyColorMap,
			layout: DefaultDateTimeLayout,
			sep:    StructSep,
			w:      os.Stdout,
			opts:   &slog.HandlerOptions{},
			mu:     &sync.Mutex{},
		},
	}
}

// WithPretty sets the `pretty`, `safe`, `color`, `layout`, and `sep` fields of the HandlerBuilder to their pretty values.
// It updates the pretty flag to true, the safe flag to true, the color map with the DefaultColorMap,
// the layout to DefaultPrettyDateTimeLayout, and the sep to PrettySep.
// Returns the updated HandlerBuilder.
func (b *HandlerBuilder) WithPretty() *HandlerBuilder {
	b.h.pretty = true
	b.h.safe = true
	b.h.color = DefaultColorMap
	b.h.layout = DefaultPrettyDateTimeLayout
	b.h.sep = PrettySep
	return b
}

// WithColor sets the color map for different log levels in the HandlerBuilder.
// If the color map is not nil, it updates the color map of the Handler.
// Returns the updated HandlerBuilder.
func (b *HandlerBuilder) WithColor(color LevelColorMap) *HandlerBuilder {
	if color != nil {
		b.h.color = color
	}
	return b
}

// WithInsecure sets the safe flag to FALSE for the HandlerBuilder.
// If the insecure flag is true, it indicates that the handler is in a safe set.
// Returns the updated HandlerBuilder.
func (b *HandlerBuilder) WithInsecure() *HandlerBuilder {
	b.h.safe = false
	return b
}

// WithTimeLayout sets the custom time layout for log messages in the HandlerBuilder.
// If the layout is not nil, it updates the time layout of the Handler.
// Returns the updated HandlerBuilder.
func (b *HandlerBuilder) WithTimeLayout(layout string) *HandlerBuilder {
	if layout != "" {
		b.h.layout = layout
	}
	return b
}

// WithSeparator sets the separator for the log message attributes in the HandlerBuilder.
// If the separator is not nil, it updates the separator of the Handler.
// Returns the updated HandlerBuilder.
func (b *HandlerBuilder) WithSeparator(sep string) *HandlerBuilder {
	if sep != "" {
		b.h.sep = sep
	}
	return b
}

// WithOptions sets the options for the HandlerBuilder.
// If the options are not nil, it updates the options of the Handler.
// Returns the updated HandlerBuilder.
func (b *HandlerBuilder) WithOptions(opts *slog.HandlerOptions) *HandlerBuilder {
	if opts != nil {
		b.h.opts = opts
	}
	return b
}

// WithWriter sets the writer for the HandlerBuilder.
// If the writer is not nil, it updates the writer of the Handler.
// Returns the updated HandlerBuilder.
func (b *HandlerBuilder) WithWriter(w io.Writer) *HandlerBuilder {
	if w != nil {
		b.h.w = w
	}
	return b
}

// WithJSON sets the JSON flag to TRUE for the HandlerBuilder.
// If the JSON flag is true, it indicates that the log messages should be formatted in JSON.
// Returns the updated HandlerBuilder.
//
// !!!Warning!!! JSON mode disable all addition otris functions. !!!Warning!!!
func (b *HandlerBuilder) WithJSON() *HandlerBuilder {
	b.h.json = true
	return b
}

// Build returns the final built Handler instance from the HandlerBuilder.
// It simply returns the value of the h field in the HandlerBuilder.
// If pretty is true, then insecure is enabled.
// If json is true, then pretty, insecure, color is disabled and sep is ','.
// Returns the final built Handler instance.
func (b *HandlerBuilder) Build() *Handler {
	if b.h.json {
		b.h.pretty = false
		b.h.safe = true
		b.h.sep = JSONSep
		b.h.color = EmptyColorMap
	}
	if b.h.pretty {
		b.h.safe = false
	}
	return b.h
}
