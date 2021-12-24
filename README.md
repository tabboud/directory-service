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
Run the Twirp server:
```sh
go run cmd/directory-service-grpc/main.go -addr localhost:8080
```

Query the server:
- Using `dsctl`:
```sh
go run cmd/dsctl/main.go -type twirp
```

- Using curl:
```sh
curl -XPOST http://localhost:8080/twirp/com.abboudlab.directoryservice.auth.AuthServiceV1/Login \
    -H 'Content-Type: application/json' \
    -d '{"username":"john","password":"super-strong"}'
```

### directory-service-grpc
Run the grpc server:
```sh
go run cmd/directory-service-grpc/main.go -addr localhost:8080 -username john -password doe
```

Query the server:
- Using `dsctl`:
```sh
go run cmd/dsctl/main.go -type grpc
```

### directory-service-conjure
```sh
go run cmd/directory-service-conjure -addr localhost:8080
```

Query the server:
- Using `dsctl`:
```sh
go run cmd/dsctl/main.go -type conjure
```

- Using curl:
```sh
curl -XPOST http://localhost:8080/v1/auth/login -k \
    -H 'Content-Type: application/json' \
    -d '{"username":"jonn","password":"super-strong"}'
```
