package logs

import (
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	stdlog "log"
	"os"
	"strings"
	"time"
)

func NoLevel(logger zerolog.Logger, level zerolog.Level) zerolog.Logger {
	return logger.Hook(NewNoLevelHook(logger.GetLevel(), level))
}

type NoLevelHook struct {
	minLevel zerolog.Level
	level    zerolog.Level
}

func NewNoLevelHook(minLevel zerolog.Level, level zerolog.Level) *NoLevelHook {
	return &NoLevelHook{minLevel: minLevel, level: level}
}

func (n NoLevelHook) Run(e *zerolog.Event, level zerolog.Level, _ string) {
	if n.minLevel > n.level {
		e.Discard()
		return
	}

	if level == zerolog.NoLevel {
		e.Str("level", n.level.String())
	}
}

func getLogWriter(filePath string) (w io.Writer) {
	w = os.Stderr

	if filePath == "stderr" {
		return
	}

	if filePath == "stdout" {
		w = zerolog.ConsoleWriter{
			Out:        w,
			TimeFormat: time.RFC3339,
			NoColor:    false,
		}
		return
	} else {
		_, _ = os.Create(filePath)

		w = &lumberjack.Logger{
			Filename: filePath,
			Compress: true,
		}

		return
	}

}

func SetupLogger(levelStr, filePath string) {
	// configure log format
	w := getLogWriter(filePath)

	if levelStr == "" {
		levelStr = "error"
	}

	logLevel, err := zerolog.ParseLevel(strings.ToLower(levelStr))
	if err != nil {
		log.Error().Err(err).
			Str("logLevel", levelStr).
			Msg("Unspecified or invalid log level, setting the level to default (ERROR)...")

		logLevel = zerolog.ErrorLevel
	}

	// create logger
	logCtx := zerolog.New(w).With().Timestamp()
	if logLevel <= zerolog.DebugLevel {
		logCtx = logCtx.Caller()
	}

	log.Logger = logCtx.Logger().Level(logLevel)
	zerolog.DefaultContextLogger = &log.Logger
	zerolog.SetGlobalLevel(logLevel)

	// configure default standard log.
	stdlog.SetFlags(stdlog.Lshortfile | stdlog.LstdFlags)
	stdlog.SetOutput(NoLevel(log.Logger, zerolog.DebugLevel))
}
