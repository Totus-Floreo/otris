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

// GetLevelName takes a slog.Level as input and returns the corresponding name as a string.
// If the level is LevelFx, it returns "FX". If the level is LevelFxError, it returns "FXError".
// For any other level, it returns the string representation of the level.
func GetLevelName(level slog.Level) (name string) {
	switch level {
	case LevelFx:
		name = "FX"
	case LevelFxError:
		name = "FXError"
	default:
		name = level.String()
	}
	return name
}

// GetColor returns the integer value of the LogColor associated with the given slog.Level in the LevelColorMap.
// If the given slog.Level is not found in the LevelColorMap, it returns the integer value of color.FgWhite (37).
func GetColor(m LevelColorMap, lvl slog.Level) int {
	color, ok := m[lvl]
	if !ok {
		// color.FgWhite == 37
		return 37
	}
	return int(color)
}
