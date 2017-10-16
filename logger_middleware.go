// Created by davidterranova on 16/10/2017.

package hserver

import (
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/davidterranova/hserver"
)

func httpLoggerMiddleware() negroni.HandlerFunc {
	return negroni.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			start := time.Now()

			rwsc := hserver.NewStatusCodeResponseWriter(rw)
			next(rwsc, r)

			logrus.WithFields(
				logrus.Fields{
					"method":        r.Method,
					"request_uri":   r.RequestURI,
					"status_code":   rwsc.Status(),
					"response_time": time.Since(start),
				},
			).Info()
		},
	)
}
