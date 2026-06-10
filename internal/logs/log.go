package logs

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func Logger(tags ...string) zerolog.Logger {
	if len(tags) == 0 {
		return log.Logger
	}

	return log.With().Str("sys", tags[0]).Logger()
}

func Verbose(on bool) {
	lvl := zerolog.InfoLevel
	if on {
		lvl = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(lvl)
}
