package server_to_server_hooks

import (
	"context"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/twitchtv/twirp"
)

const (
	missingJWT              = "Missing JWT"
	missingPublicKey        = "Missing public key"
	notAServerToServerToken = "Not a server-to-server token"
)

func Auth() *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			err := serviceToServiceJWTIsValid(ctx)
			if err != nil {
				return context.Background(), twirp.NewError(twirp.Unauthenticated, err.Error())
			}
			return ctx, nil
		},
	}
}

func Logging() *twirp.ServerHooks {
	return &twirp.ServerHooks{
		ResponseSent: func(ctx context.Context) {
			methodName, _ := twirp.MethodName(ctx)
			requestParams, _ := ctx.Value("requestBody").(map[string]interface{})
			logrus.WithFields(logrus.Fields{"params": requestParams, "path": methodName}).Debug()
		},
		Error: func(ctx context.Context, error twirp.Error) context.Context {
			methodName, _ := twirp.MethodName(ctx)
			requestParams, _ := ctx.Value("requestBody").(map[string]interface{})
			logrus.WithFields(logrus.Fields{"params": requestParams, "hook": "Logging:Error",
				"path": methodName}).Error(error)

			return ctx
		},
	}
}

func serviceToServiceJWTIsValid(ctx context.Context) error {
	token, ok := ctx.Value("jwtToken").(string)
	if !ok {
		return errors.New(missingJWT)
	}
	token = strings.Split(token, "Bearer ")[1]
	pubKey, err := ioutil.ReadFile("/Users/charlie/jwt/jwtRS256.key.pub")
	if err != nil {
		return errors.New(missingPublicKey)
	}
	jwtToken := NewJWT(pubKey)
	_, err = jwtToken.Validate(token)
	if err != nil {
		return err
	}
	return nil
}
