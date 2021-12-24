# Directory Service

A simple service used to experiment with [twirp](https://github.com/twitchtv/twirp).

## Development

Run the following command to re-generate the server/client code:
```sh
protoc --go_out=paths=source_relative:. --twirp_out=paths=source_relative:. rpc/authservice/service.proto
```
