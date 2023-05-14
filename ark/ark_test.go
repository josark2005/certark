package ark_test

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestLog(T *testing.T) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("Hello World")
}
