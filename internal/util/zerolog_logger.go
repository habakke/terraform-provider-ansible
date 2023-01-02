package util

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
)

type ZerologLogger struct {
}

func NewZerologLogger() ZerologLogger {
	return ZerologLogger{}
}

func (t ZerologLogger) Tracef(ctx context.Context, msg string, additionalFields ...interface{}) {
	log.Trace().Msgf(fmt.Sprintf(msg, additionalFields...))
}

func (t ZerologLogger) Debugf(ctx context.Context, msg string, additionalFields ...interface{}) {
	log.Debug().Msgf(fmt.Sprintf(msg, additionalFields...))
}

func (t ZerologLogger) Warnf(ctx context.Context, msg string, additionalFields ...interface{}) {
	log.Warn().Msgf(fmt.Sprintf(msg, additionalFields...))
}

func (t ZerologLogger) Errorf(ctx context.Context, msg string, additionalFields ...interface{}) {
	log.Error().Msgf(fmt.Sprintf(msg, additionalFields...))
}

func (t ZerologLogger) Infof(ctx context.Context, msg string, additionalFields ...interface{}) {
	log.Info().Msgf(fmt.Sprintf(msg, additionalFields...))
}
