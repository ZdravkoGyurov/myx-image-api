package middlewares

import (
	"net/http"

	"github.com/ZdravkoGyurov/myx-image-api/pkg/api/response"
	"github.com/ZdravkoGyurov/myx-image-api/pkg/errors"
)

func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			panicErr := recover()
			if panicErr != nil {
				err := errors.Newf("recovered from panic: %s", panicErr)
				response.SendError(writer, request, err)
			}
		}()

		next.ServeHTTP(writer, request)
	})
}
