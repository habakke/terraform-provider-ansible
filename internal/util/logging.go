package util

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"time"
)

//nolint:deadcode,may be used later
//lint:file-ignore U1000 May be used later
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func ConfigureLogging(logLevel string, logCaller bool) {
	l, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		l = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(l)
	if logCaller {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Caller().Logger()
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}
}

// ConfigureTerraformProviderLogging configures zerolog according to the format expected by terraform
// according to https://www.terraform.io/docs/extend/debugging.html#log-based-debugging
func ConfigureTerraformProviderLogging(logCaller bool) {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	output := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339, NoColor: true, PartsExclude: []string{zerolog.TimestampFieldName}}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("[%s]", i))
	}
	if logCaller {
		log.Logger = log.Output(output).With().Caller().Logger()
	} else {
		log.Logger = log.Output(output)
	}
}
