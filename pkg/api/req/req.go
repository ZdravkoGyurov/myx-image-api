package req

import (
	"context"
	"net/http"
)

type correlationIDKey struct{}

var (
	CorrelationIDKey correlationIDKey
)

func AddCorrelationID(r *http.Request, correlationID string) {
	*r = *r.WithContext(context.WithValue(r.Context(), CorrelationIDKey, correlationID))
}

func GetCorrelationID(r *http.Request) (string, bool) {
	correlationID, ok := r.Context().Value(CorrelationIDKey).(string)
	return correlationID, ok
}
