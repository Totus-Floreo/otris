package otris

import "log/slog"

// LevelFx and LevelFxError represents a custom logging level for FxHandler. It has a value of -8 and -7.
const (
	LevelFx      = slog.Level(-8)
	LevelFxError = slog.Level(-7)
	LevelDebug   = slog.LevelDebug
	LevelInfo    = slog.LevelInfo
	LevelWarning = slog.LevelWarn
	LevelError   = slog.LevelError
)
