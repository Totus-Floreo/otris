package otris

import (
	"github.com/fatih/color"
	"log/slog"
)

// LogColor defines a single SGR Code
type LogColor int

type LevelColorMap map[slog.Level]LogColor

// DefaultColorMap is the default color mapping used for logging.
var DefaultColorMap = LevelColorMap{
	LevelFx:      LogColor(color.FgCyan),
	LevelFxError: LogColor(color.FgHiRed),
	LevelDebug:   LogColor(color.FgBlue),
	LevelInfo:    LogColor(color.FgHiGreen),
	LevelWarning: LogColor(color.FgYellow),
	LevelError:   LogColor(color.FgRed),
}

// EmptyColorMap is the empty color mapping used for safe logging.
var EmptyColorMap = LevelColorMap{}

//TODO WIP in v2 Coloring value in logs
/*
// LogKey represents a key used for logging.
type LogKey string

// LogValue represents a value used for logging.
type LogValue any

// ColorMap represents a mapping of LogKey to a mapping of LogValue to LogColor.
// It is used to define the color scheme for different log keys and values.
type ColorMapV2 map[LogKey]map[LogValue]LogColor

// DefaultColorMap is the default color mapping used for logging.
// TODO Implement SQL and HTTP error codes in v2
var DefaultColorMapV2 = ColorMapV2{
	LogKey(slog.LevelKey): {
		LogValue(LevelFx):      LogColor(color.FgCyan),
		LogValue(LevelFxError): LogColor(color.FgHiRed),
		LogValue(LevelDebug):   LogColor(color.FgHiMagenta),
		LogValue(LevelInfo):    LogColor(color.FgHiGreen),
		LogValue(LevelWarning): LogColor(color.FgHiYellow),
		LogValue(LevelError):   LogColor(color.FgRed),
	},
	LogKey("httpcode"): {
		LogValue(200): LogColor(color.FgGreen),
		LogValue(404): LogColor(color.FgYellow),
		LogValue(500): LogColor(color.FgRed),
	},
}
*/
