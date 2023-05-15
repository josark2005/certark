package ark_test

import (
	"testing"

	"github.com/jokin1999/certark/ark"
)

func TestLog(T *testing.T) {
	ark.Trace().Msg("this is trace")
	ark.Debug().Msg("this is debug")
	ark.Info().Str("key", "value").Msg("this is info")
	ark.Warn().Msg("this is warn")
	ark.Error().Msg("this is error")
}
