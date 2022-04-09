gen:
    # Auto-generate code
	protoc --go_out=. --twirp_out=. rpc/haberdasher/service.proto

upgrade:
    # Upgrade dependencies if using modules
	go get -u