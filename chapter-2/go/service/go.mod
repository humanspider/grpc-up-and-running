module productinfo/service

go 1.23.0

toolchain go1.24.1

require (
	github.com/gofrs/uuid v4.4.0+incompatible
	google.golang.org/grpc v1.71.1
	google.golang.org/protobuf v1.36.6
)

require (
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250407143221-ac9807e6c755 // indirect
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.5.1 // indirect
)

tool (
	google.golang.org/grpc/cmd/protoc-gen-go-grpc
	google.golang.org/protobuf/cmd/protoc-gen-go
)
