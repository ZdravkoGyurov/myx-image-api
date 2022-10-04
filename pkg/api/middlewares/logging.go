package middlewares

import (
	"net/http"

	"github.com/ZdravkoGyurov/myx-image-api/pkg/api/req"
	"github.com/ZdravkoGyurov/myx-image-api/pkg/log"

	"github.com/google/uuid"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.DefaultLogger().With().Logger()
		loggerContext := logger.WithContext(r.Context())
		*r = *r.WithContext(loggerContext)
		next.ServeHTTP(w, r)
	})
}

func CorrelationIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.RequestLogger(r)

		correlationID := r.Header.Get(log.CorrelationIDHeader)
		if _, err := uuid.Parse(correlationID); err != nil {
			if correlationID != "" {
				logger.Warn().Msgf("Invalid correlation ID in request header: %s", correlationID)
			}
			correlationID = uuid.NewString()
		}

		req.AddCorrelationID(r, correlationID)
		log.SetCorrelationID(logger, correlationID)

		w.Header().Set(log.CorrelationIDHeader, correlationID)
		next.ServeHTTP(w, r)
	})
}
