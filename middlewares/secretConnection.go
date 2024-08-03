package middlewares

import (
	"net/http"

	"github.com/ynoacamino/ynoa-uploader/config"
)

func ConncetionSecret(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := r.Header.Get("Secret-Connection")
		if secret != config.SECRET_CONNECTION {
			http.Error(w, "Not Authorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
