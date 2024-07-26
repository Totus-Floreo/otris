package otris

import (
	"bytes"
	"context"
	"github.com/fatih/color"
	"log/slog"
	"testing"
	"time"
)

func TestNewStructHandler(t *testing.T) {
	ctx := context.Background()
	groupKey := "groupKey"
	groupValue := "Group1"
	timeKey := "timestamp"
	timeValue := time.Now()

	preAttrs := []slog.Attr{slog.Int("pre", 0), slog.String(groupKey, groupValue), slog.Time(timeKey, timeValue)}
	attrs := []slog.Attr{slog.Int("a", 1), slog.String("b", "two"), slog.String(groupKey, groupValue), slog.Time(timeKey, timeValue)}
	emptyAttrs := []slog.Attr{}

	handlerOpts := []slog.HandlerOptions{
		{Level: slog.LevelDebug},
		{Level: slog.LevelInfo},
		{Level: slog.LevelWarn},
		{Level: slog.LevelError},
	}

	// Test cases
	cases := []struct {
		name  string
		attrs []slog.Attr
		opts  slog.HandlerOptions
		want  string
	}{
		{
			name:  "Case 1",
			attrs: preAttrs,
			opts:  handlerOpts[0],
			want:  "",
		},
		{
			name:  "Case 2",
			attrs: attrs,
			opts:  handlerOpts[1],
			want:  "",
		},
		{
			name:  "Case 3",
			attrs: emptyAttrs,
			opts:  handlerOpts[2],
			want:  "",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			var got1 bytes.Buffer
			var got2 bytes.Buffer
			var h1 slog.Handler = NewStructHandler(&got1, &test.opts)
			var h2 slog.Handler = slog.NewTextHandler(&got2, &test.opts)

			r := slog.NewRecord(time.Time{}, LevelInfo, "message", 0)
			r.AddAttrs(test.attrs...)

			if err := h1.Handle(ctx, r); err != nil {
				t.Fatal(err)
			}
			if err := h2.Handle(ctx, r); err != nil {
				t.Fatal(err)
			}

			if got1.String() != got2.String() {
				t.Errorf("\ngot  %s\nwant %s", got1.String(), got2.String())
			}
		})
	}
}

func TestNewJSONHandler(t *testing.T) {
	ctx := context.Background()
	groupKey := "groupKey"
	groupValue := "Group1"
	timeKey := "timestamp"
	timeValue := time.Now()

	preAttrs := []slog.Attr{slog.Int("pre", 0), slog.String(groupKey, groupValue), slog.Time(timeKey, timeValue)}
	attrs := []slog.Attr{slog.Int("a", 1), slog.String("b", "two"), slog.String(groupKey, groupValue), slog.Time(timeKey, timeValue)}
	emptyAttrs := []slog.Attr{}

	handlerOpts := []slog.HandlerOptions{
		{Level: slog.LevelDebug},
		{Level: slog.LevelInfo},
		{Level: slog.LevelWarn},
		{Level: slog.LevelError},
	}

	// Test cases
	cases := []struct {
		name  string
		attrs []slog.Attr
		opts  slog.HandlerOptions
		want  string
	}{
		{
			name:  "Case 1",
			attrs: preAttrs,
			opts:  handlerOpts[0],
			want:  "",
		},
		{
			name:  "Case 2",
			attrs: attrs,
			opts:  handlerOpts[1],
			want:  "",
		},
		{
			name:  "Case 3",
			attrs: emptyAttrs,
			opts:  handlerOpts[2],
			want:  "",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			var got1 bytes.Buffer
			var got2 bytes.Buffer
			var h1 slog.Handler = NewJSONHandler(&got1, &test.opts)
			var h2 slog.Handler = slog.NewJSONHandler(&got2, &test.opts)

			r := slog.NewRecord(time.Time{}, LevelInfo, "message", 0)
			r.AddAttrs(test.attrs...)

			if err := h1.Handle(ctx, r); err != nil {
				t.Fatal(err)
			}
			if err := h2.Handle(ctx, r); err != nil {
				t.Fatal(err)
			}

			if got1.String() != got2.String() {
				t.Errorf("\ngot  %s\nwant %s", got1.String(), got2.String())
			}
		})
	}
}

func TestHandlerBuilder(t *testing.T) {
	ctx := context.Background()
	groupKey := "groupKey"
	groupValue := "Group1"
	timeKey := "timestamp"
	timeValue := time.Now()

	preAttrs := []slog.Attr{slog.Int("pre", 0), slog.String(groupKey, groupValue), slog.Time(timeKey, timeValue)}
	attrs := []slog.Attr{slog.Int("a", 1), slog.String("b", "two"), slog.String(groupKey, groupValue), slog.Time(timeKey, timeValue)}
	emptyAttrs := []slog.Attr{}

	colorMap := LevelColorMap{
		slog.LevelDebug: LogColor(color.FgWhite),
		slog.LevelInfo:  LogColor(color.FgGreen),
		slog.LevelWarn:  LogColor(color.FgYellow),
		slog.LevelError: LogColor(color.FgRed),
	}

	// Test case
	cases := []struct {
		name     string
		attrs    []slog.Attr
		colormap LevelColorMap
	}{
		{
			name:     "Case 1",
			attrs:    preAttrs,
			colormap: DefaultColorMap,
		},
		{
			name:     "Case 2",
			attrs:    attrs,
			colormap: colorMap,
		},
		{
			name:     "Case 3",
			attrs:    emptyAttrs,
			colormap: colorMap,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			var got1 bytes.Buffer
			var got2 bytes.Buffer
			var got3 bytes.Buffer

			var h1 slog.Handler = NewHandler(&got1, test.colormap, false, "", "", nil)
			var h2 slog.Handler = NewHandlerBuilder().WithPretty().WithWriter(&got2).WithColor(test.colormap).Build()
			var h3 slog.Handler = NewPrettyHandler(&got3, nil)

			r := slog.NewRecord(time.Now(), LevelInfo, "message", 0)
			r.AddAttrs(test.attrs...)

			if err := h1.Handle(ctx, r); err != nil {
				t.Fatal(err)
			}
			if err := h2.Handle(ctx, r); err != nil {
				t.Fatal(err)
			}
			if err := h3.Handle(ctx, r); err != nil {
				t.Fatal(err)
			}

			if got1.String() != got2.String() {
				t.Errorf("\ngot  %s\nwant %s", got1.String(), got2.String())
			}

			if got2.String() != got3.String() {
				t.Errorf("\ngot  %s\nwant %s", got2.String(), got3.String())
			}
		})
	}
}

func TestNewHandlerBuilderWithJSON(t *testing.T) {
	ctx := context.Background()
	groupKey := "groupKey"
	groupValue := "Group1"
	timeKey := "timestamp"
	timeValue := time.Now()

	preAttrs := []slog.Attr{slog.Int("pre", 0), slog.String(groupKey, groupValue), slog.Time(timeKey, timeValue)}
	attrs := []slog.Attr{slog.Int("a", 1), slog.String("b", "two"), slog.String(groupKey, groupValue), slog.Time(timeKey, timeValue)}
	emptyAttrs := []slog.Attr{}

	handlerOpts := []slog.HandlerOptions{
		{Level: slog.LevelDebug},
		{Level: slog.LevelInfo},
		{Level: slog.LevelWarn},
		{Level: slog.LevelError},
	}

	colorMap := LevelColorMap{
		slog.LevelDebug: LogColor(color.FgWhite),
		slog.LevelInfo:  LogColor(color.FgGreen),
		slog.LevelWarn:  LogColor(color.FgYellow),
		slog.LevelError: LogColor(color.FgRed),
	}

	// Test cases
	cases := []struct {
		name  string
		attrs []slog.Attr
		opts  slog.HandlerOptions
		want  string
	}{
		{
			name:  "Case 1",
			attrs: preAttrs,
			opts:  handlerOpts[0],
			want:  "",
		},
		{
			name:  "Case 2",
			attrs: attrs,
			opts:  handlerOpts[1],
			want:  "",
		},
		{
			name:  "Case 3",
			attrs: emptyAttrs,
			opts:  handlerOpts[2],
			want:  "",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			var got1 bytes.Buffer
			var got2 bytes.Buffer
			var h1 slog.Handler = NewHandlerBuilder().WithWriter(&got1).WithColor(colorMap).WithJSON().Build()
			var h2 slog.Handler = slog.NewJSONHandler(&got2, &test.opts)

			r := slog.NewRecord(time.Time{}, LevelInfo, "message", 0)
			r.AddAttrs(test.attrs...)

			if err := h1.Handle(ctx, r); err != nil {
				t.Fatal(err)
			}
			if err := h2.Handle(ctx, r); err != nil {
				t.Fatal(err)
			}

			if got1.String() != got2.String() {
				t.Errorf("\ngot  %s\nwant %s", got1.String(), got2.String())
			}
		})
	}
}

func TestNewPrettyHandler(t *testing.T) {
	ctx := context.Background()
	groupKey := "groupKey"
	groupValue := "Group1"
	timeKey := "timestamp"
	timeValue := time.Now()

	preAttrs := []slog.Attr{slog.Int("pre", 0), slog.String(groupKey, groupValue), slog.Time(timeKey, timeValue)}
	attrs := []slog.Attr{slog.Int("a", 1), slog.String("b", "two"), slog.String(groupKey, groupValue), slog.Time(timeKey, timeValue)}
	emptyAttrs := []slog.Attr{}

	// Test case
	cases := []struct {
		name  string
		attrs []slog.Attr
	}{
		{
			name:  "Case 1",
			attrs: preAttrs,
		},
		{
			name:  "Case 2",
			attrs: attrs,
		},
		{
			name:  "Case 3",
			attrs: emptyAttrs,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			var got bytes.Buffer

			var h1 slog.Handler = NewPrettyHandler(&got, nil)

			r := slog.NewRecord(time.Now(), LevelInfo, "message", 0)
			r.AddAttrs(test.attrs...)

			if err := h1.Handle(ctx, r); err != nil {
				t.Fatal(err)
			}

			t.Logf("The log looks like: \n%s", got.String())
		})
	}
}
