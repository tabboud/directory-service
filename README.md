# Directory Service

A simple service used to experiment with different RPC frameworks:
- [twirp](https://github.com/twitchtv/twirp)
- [gRPC](https://github.com/grpc/grpc-go)
- [conjure](https://github.com/palantir/conjure-go)

## Development

Re-generate the proto defined server/client code:
```sh
protoc --go_out=paths=source_relative:. --go-grpc_out=. --go-grpc_opt=paths=source_relative --twirp_out=paths=source_relative:. rpc/authservice
```

Re-generate conjure code:
```sh
./godelw conjure
```

### directory-service-twip
```
Run the Twirp server:
```sh
go run cmd/directory-service-grpc/main.go -addr localhost:8080
```

Use curl to query the server:
```sh
curl -XPOST http://localhost:8080/twirp/com.abboudlab.directoryservice.auth.AuthServiceV1/Login \
    -H 'Content-Type: application/json' \
    -d '{"username":"gooduser","password":"test-password"}'
```

### directory-service-grpc
Run the grpc server:
```sh
go run cmd/directory-service-grpc/main.go -addr localhost:8080 -username john -password doe
```

Use the dsctl cli to query the grpc server:
```sh
go run cmd/dsctl/main.go -addr localhost:8080
```
