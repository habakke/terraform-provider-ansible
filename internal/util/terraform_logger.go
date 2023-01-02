package util

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type TerraformLogger struct {
}

func NewTerraformLogger() TerraformLogger {
	return TerraformLogger{}
}

func (t TerraformLogger) Tracef(ctx context.Context, msg string, additionalFields ...interface{}) {
	tflog.Trace(ctx, fmt.Sprintf(msg, additionalFields...))
}

func (t TerraformLogger) Debugf(ctx context.Context, msg string, additionalFields ...interface{}) {
	tflog.Debug(ctx, fmt.Sprintf(msg, additionalFields...))
}

func (t TerraformLogger) Warnf(ctx context.Context, msg string, additionalFields ...interface{}) {
	tflog.Warn(ctx, fmt.Sprintf(msg, additionalFields...))
}

func (t TerraformLogger) Errorf(ctx context.Context, msg string, additionalFields ...interface{}) {
	tflog.Error(ctx, fmt.Sprintf(msg, additionalFields...))
}

func (t TerraformLogger) Infof(ctx context.Context, msg string, additionalFields ...interface{}) {
	tflog.Info(ctx, fmt.Sprintf(msg, additionalFields...))
}
