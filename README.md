#Twirp Best Practices

##JWT

* Use standard JWT claims, and pass user_id or admin_id in the token

###user
* Claims
  * exp: expiration (unix time)
  * iat: issued at timestamp (unix time)
  * iss: issuer ("charlie.com")
  * aud: audience ("user")
  * user_id: user_id (integer)

###admin
* Claims
  * exp: expiration (unix time)
  * iat: issued at timestamp (unix time)
  * iss: issuer ("charlie.com")
  * aud: audience ("admin")
  * admin_id: admin_id (integer)

###server-to-server
* Claims
  * exp: expiration (unix time)
  * iat: issued at timestamp (unix time)
  * iss: issuer ("charlie.com")
  * aud: audience ("server-to-server")

###ChainHooks
Create as many different `ChainHook` objects as you need. Here we create different hook objects for user, admin, and server-to-server

**cmd/server/main.go**

```
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
```
Each `ChainHook` object has access to these parts of the request / response lifecycle:

**twirp/server_options.go**
```
type ServerHooks struct {
	// RequestReceived is called as soon as a request enters the Twirp
	// server at the earliest available moment.
	RequestReceived func(context.Context) (context.Context, error)

	// RequestRouted is called when a request has been routed to a
	// particular method of the Twirp server.
	RequestRouted func(context.Context) (context.Context, error)

	// ResponsePrepared is called when a request has been handled and a
	// response is ready to be sent to the client.
	ResponsePrepared func(context.Context) context.Context

	// ResponseSent is called when all bytes of a response (including an error
	// response) have been written. Because the ResponseSent hook is terminal, it
	// does not return a context.
	ResponseSent func(context.Context)

	// Error hook is called when an error occurs while handling a request. The
	// Error is passed as argument to the hook.
	Error func(context.Context, Error) context.Context
}
```

Put it all together: Share twirp service code and handlers if you like, but use distinct paths + ChainHooks for user, admin, server-to-server

**cmd/server/main.go**

```
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
```