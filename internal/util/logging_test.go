package util

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"testing"
)

func TestTerraformProviderLoggingFormat(t *testing.T) {
	ConfigureTerraformProviderLogging("trace", true)

	log.Error().Err(fmt.Errorf("test error message")).Msg("This is an error level message")
	log.Info().Msg("This is an info level message")
	log.Warn().Msg("This is an warn level message")
	log.Trace().Msg("This is a trace level message")
}
