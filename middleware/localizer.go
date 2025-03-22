package middleware

import (
	"fmt"
	"net/http"
	"time"

	translator "github.com/h1sashin/go-app/i18n"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func Localizer(bundle *i18n.Bundle) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()
			defer func() {
				elapsed := time.Since(now)
				fmt.Printf("Request processed in %s\n", elapsed)
			}()
			h := r.Header.Get("Accept-Language")
			loc := i18n.NewLocalizer(bundle, h)

			ctx := translator.InjectLocalizer(r.Context(), loc)

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
