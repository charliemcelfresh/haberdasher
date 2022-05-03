package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/sirupsen/logrus"
)

const (
	Authorization = "Authorization"
)

func AddRequestBodyToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logrus.Error(err)
		}
		cloneOne := ioutil.NopCloser(bytes.NewBuffer(buf))
		cloneTwo := ioutil.NopCloser(bytes.NewBuffer(buf))
		clonedBody, err := ioutil.ReadAll(cloneOne)
		if err != nil {
			logrus.Error(err)
		}
		var m map[string]interface{}
		err = json.Unmarshal(clonedBody, &m)
		if err != nil {
			logrus.Error(err)
		}
		ctx = context.WithValue(ctx, "requestBody", m)
		r.Body = cloneTwo
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func AddJwtTokenToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get(Authorization)
		if auth == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		} else {
			auth := r.Header.Get(Authorization)
			ctx := r.Context()
			ctx = context.WithValue(ctx, "jwtToken", auth)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		}
	})
}
