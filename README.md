## Sample GRPC Client/Server

Sample GRPC client/server with protobuf

### Setup
- Add protoc to PATH and unzip: https://github.com/protocolbuffers/protobuf/releases/tag/v3.10.1
- go get -u google.golang.org/grpc
- go get -u github.com/golang/protobuf/protoc-gen-go

### Run
`go run server/main.go`

`go run client/main.go`

`localhost:8080/operations/add/2/3`

`localhost:8080/operations/add/2/3`
