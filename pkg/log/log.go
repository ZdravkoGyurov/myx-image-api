package log

import (
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	CorrelationIDHeader    = "X-Correlation-Id"
	correlationIDAttribute = "correlationId"
	componentAttribute     = "component"
	component              = "grader"
)

var logger zerolog.Logger

func init() {
	var loggerOutput io.Writer
	if os.Getenv("LOCAL_DEV") == "true" {
		loggerOutput = zerolog.NewConsoleWriter()
	} else {
		loggerOutput = os.Stdout
	}
	logger = zerolog.New(loggerOutput).With().Timestamp().Str(componentAttribute, component).Logger()
}

func SetCorrelationID(logger *zerolog.Logger, correlationID string) {
	logger.UpdateContext(func(logCtx zerolog.Context) zerolog.Context {
		return logCtx.Str(correlationIDAttribute, correlationID)
	})
}

func DefaultLogger() *zerolog.Logger {
	return &logger
}

func RequestLogger(r *http.Request) *zerolog.Logger {
	return log.Ctx(r.Context())
}

var CtxLogger = log.Ctx
