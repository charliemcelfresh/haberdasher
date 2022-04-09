package main

import (
	"charliemcelfresh/haberdasher/internal/haberdasherserver"
	"charliemcelfresh/haberdasher/rpc/haberdasher"
	"context"
	"log"
	"net/http"

	"github.com/twitchtv/twirp"
)

func main() {

	server := &haberdasherserver.Server{} // implements Haberdasher interface

	mux := http.NewServeMux()

	// http(s)://<host>:/v1/user/haberdasher.Haberdasher/MakeHat
	// http(s)://<host>:/v1/user/haberdasher.Haberdasher/HelloWorld
	userHandler := haberdasher.NewHaberdasherServer(server, twirp.WithServerPathPrefix("/v1/users"), NewUserServerHooks())
	mux.Handle(userHandler.PathPrefix(), userHandler)

	// http(s)://<host>:/v1/admin/haberdasher.Haberdasher/MakeHat
	// http(s)://<host>:/v1/admin/haberdasher.Haberdasher/HelloWorld
	adminHandler := haberdasher.NewHaberdasherServer(server, twirp.WithServerPathPrefix("/v1/admin"), NewAdminServerHooks())
	mux.Handle(adminHandler.PathPrefix(), adminHandler)

	// http(s)://<host>:/v1/internal/haberdasher.Haberdasher/MakeHat
	// http(s)://<host>:/v1/internal/haberdasher.Haberdasher/HelloWorld
	serviceToServiceHandler := haberdasher.NewHaberdasherServer(server, twirp.WithServerPathPrefix("/v1/internal"),
		NewServiceToServiceServerHooks())
	mux.Handle(serviceToServiceHandler.PathPrefix(), serviceToServiceHandler)

	http.ListenAndServe(":8080", mux)
}

/*
	Separate hooks for User, Admin and Service-to-service, for auth, user_audit_trail, events, and logging

	Each knows how to handle its own auth, and rejects before the handler functions are entered -- keep
	auth code out of handlers, and pass needed information (userID, adminID, calling service name, etc.)
	along in the context
*/
func NewUserServerHooks() *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestRouted: func(ctx context.Context) (context.Context, error) {
			log.Println("User authenticated")
			return ctx, nil
		},
	}
}

func NewAdminServerHooks() *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestRouted: func(ctx context.Context) (context.Context, error) {
			log.Println("Admin authenticated")
			return ctx, nil
		},
	}
}

func NewServiceToServiceServerHooks() *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestRouted: func(ctx context.Context) (context.Context, error) {
			log.Println("Internal authenticated")
			return ctx, nil
		},
	}
}
