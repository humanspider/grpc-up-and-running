# Set Up

## Install gRPC

Download the latest release zip from here: https://github.com/protocolbuffers/protobuf/releases/

Uncompress the file `protobuf-<version>-<arch>.<tar.gz or zip>`.

Place the extracted contents of `bin` in `usr/local/bin` and `include` in `usr/local/include`.

## Golang service

In your service project directory, create a new Go module: `go mod init productinfo/service`

### Install gRPC library

Run `go get -u google.golang.org/grpc`

### Install protoc Go plugin

Golang protoc plugin has migrated to a newer version, called Opaque API. 

Run `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`and `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc` to install the plugin in `$GOBIN`(defaults to $GOPATH/bin).

### Generate Go output

Navigate to the the root of the current chapter directory.

Run the compilation command: `protoc --proto_path=proto --go_out=service --go_opt=module=productinfo/service product_info.proto --go-grpc_out=service --go-grpc_opt=module=productinfo/service`. This will output the generated Go protobuf inside of service/ecommerce.

Here's a breakdown of the command
1. `--proto_path` is the directory of the `.proto` files
2. `--go_out=service` specifies the service directory as the root for the output
3. By default, the behavior is to place the output file in the directory named after the Go packages import path, such as one provided by the `go_package` option within the `.proto` file (can also be specified with the --go_opt=paths=import option). `--go_opt=module=productinfo/service`follows that logic, but removes the specified module from the output filename (`productinfo/service/ecommerce` becomes `ecommerce`).
4. The go-grpc flags serve the same purpose as the other flags, but for the grpc client and server output

The current project layout places the protobuf definitions under the chapter-x/proto directory. If you wished to keep the protobuf definitions in the service project itself (`service/proto`), then you would use this command from the `service` root: `protoc --proto_path=proto --go_out=. --go_opt=module=productinfo/service`.

### Proto 2023 message changes

Protobuf Edition 2023 introduces a more explicit and consistent way to handle optional fields. Before Edition 2023, unset optional fields would be represented by their zero value. This would cause issues when the zero value couldn't be determined to be optional by the message recipient. Edition 2023 represents values through pointers, so that their existance can be verified using a nil check. To set these pointer values, use the `google.golang.org/protobuf/proto` scalar constructors to create a new value and return a pointer.

## Java service

