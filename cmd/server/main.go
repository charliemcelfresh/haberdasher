package main

import (
	"charliemcelfresh/haberdasher/internal/admin_hooks"
	"charliemcelfresh/haberdasher/internal/haberdasherserver"
	"charliemcelfresh/haberdasher/internal/middlewares"
	"charliemcelfresh/haberdasher/internal/server_to_server_hooks"
	"charliemcelfresh/haberdasher/internal/user_hooks"
	"charliemcelfresh/haberdasher/rpc/haberdasher"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/twitchtv/twirp"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {

	server := &haberdasherserver.Server{} // implements Haberdasher interface

	mux := http.NewServeMux()

	userChainHooks := twirp.ChainHooks(
		user_hooks.Auth(),
		user_hooks.Logging(),
	)

	serverToServerChainHooks := twirp.ChainHooks(
		server_to_server_hooks.Auth(),
		server_to_server_hooks.Logging(),
	)

	adminChainHooks := twirp.ChainHooks(
		admin_hooks.Auth(),
		admin_hooks.Audit(),
		admin_hooks.Logging(),
	)

	// http(s)://<host>:/v1/user/haberdasher.Haberdasher/MakeHat
	// http(s)://<host>:/v1/user/haberdasher.Haberdasher/HelloWorld
	userHandler := haberdasher.NewHaberdasherServer(server, twirp.WithServerPathPrefix("/v1/user"), userChainHooks)
	mux.Handle(userHandler.PathPrefix(), middlewares.AddRequestBodyToContext(middlewares.AddJwtTokenToContext(
		userHandler)))

	// http(s)://<host>:/v1/admin/haberdasher.Haberdasher/MakeHat
	// http(s)://<host>:/v1/admin/haberdasher.Haberdasher/HelloWorld
	adminHandler := haberdasher.NewHaberdasherServer(server, twirp.WithServerPathPrefix("/v1/admin"), adminChainHooks)
	mux.Handle(adminHandler.PathPrefix(), middlewares.AddRequestBodyToContext(middlewares.AddJwtTokenToContext(
		adminHandler)))

	// http(s)://<host>:/v1/internal/haberdasher.Haberdasher/MakeHat
	// http(s)://<host>:/v1/internal/haberdasher.Haberdasher/HelloWorld
	serviceToServiceHandler := haberdasher.NewHaberdasherServer(server, twirp.WithServerPathPrefix("/v1/internal"),
		serverToServerChainHooks)
	mux.Handle(serviceToServiceHandler.PathPrefix(), middlewares.AddRequestBodyToContext(middlewares.AddJwtTokenToContext(
		serviceToServiceHandler)))

	http.ListenAndServe(":8080", mux)
}

/*
	Separate hooks for User, Admin and Service-to-service, for auth, user_audit_trail, events, and logging

	Each knows how to handle its own auth, and rejects before the handler functions are entered -- keep
	auth code out of handlers, and pass needed information (userID, adminID, calling service name, etc.)
	along in the context
*/
