package admin_hooks

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/twitchtv/twirp"
)

const (
	missingJWT       = "Missing JWT"
	missingPublicKey = "Missing public key"
)

func Auth() *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			claims, err := adminJWTIsValid(ctx)
			if err != nil {
				return context.Background(), twirp.NewError(twirp.Unauthenticated, err.Error())
			}
			err = validateAdmin(claims)
			if err != nil {
				return context.Background(), twirp.NewError(twirp.Unauthenticated, err.Error())
			}
			return ctx, nil
		},
	}
}

func Audit() *twirp.ServerHooks {
	return &twirp.ServerHooks{
		ResponseSent: func(ctx context.Context) {
			methodName, _ := twirp.MethodName(ctx)
			msg := fmt.Sprintf("Admin performed %s", methodName)
			logrus.Debug(msg)
		},
	}
}

func Logging() *twirp.ServerHooks {
	return &twirp.ServerHooks{
		ResponseSent: func(ctx context.Context) {
			methodName, _ := twirp.MethodName(ctx)
			requestParams, _ := ctx.Value("requestBody").(map[string]interface{})
			logrus.WithFields(logrus.Fields{"params": requestParams,
				"path": methodName}).Debug()
		},
		Error: func(ctx context.Context, error twirp.Error) context.Context {
			methodName, _ := twirp.MethodName(ctx)
			requestParams, _ := ctx.Value("requestBody").(map[string]interface{})
			logrus.WithFields(logrus.Fields{"params": requestParams, "hook": "Logging:Error",
				"methodName": methodName}).Error(error)
			return ctx
		},
	}

}

func adminJWTIsValid(ctx context.Context) (Claims, error) {
	token, ok := ctx.Value("jwtToken").(string)
	if !ok {
		return Claims{}, errors.New(missingJWT)
	}
	token = strings.Split(token, "Bearer ")[1]
	pubKey, err := ioutil.ReadFile("/Users/charlie/jwt/jwtRS256.key.pub")
	if err != nil {
		return Claims{}, errors.New(missingPublicKey)
	}
	jwtToken := NewJWT(pubKey)
	claims, err := jwtToken.Validate(token)
	if err != nil {
		return Claims{}, err
	}
	return claims, nil
}

func validateAdmin(claims Claims) error {
	// look up user in the DB, make sure they exist without restriction
	return nil
}
