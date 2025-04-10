# Set Up

## Install protoc

Download the latest release zip from here: https://github.com/protocolbuffers/protobuf/releases/

Uncompress the file `protobuf-<version>-<arch>.<tar.gz or zip>`.

Place the extracted contents of `bin` in `usr/local/bin` and `include` in `usr/local/include`.

## Golang Generation

In your service project directory, create a new Go module: `go mod init productinfo/service`

### Install gRPC library

Run `go get -u google.golang.org/grpc`

### Install protoc Go plugin

Golang protoc plugin has migrated to a newer version, called Opaque API. 

Run `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`and `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc` to install the plugin in `$GOBIN`(defaults to $GOPATH/bin). Include `$GOBIN` in your path so that protoc can find the plugins.

### Generate Go output

Navigate to the the root of the current chapter directory.

Run the compilation command: `protoc --proto_path=proto --go_out=service --go_opt=module=productinfo/service product_info.proto --go-grpc_out=service --go-grpc_opt=module=productinfo/service`. This will output the generated Go protobuf inside of service/ecommerce.

Here's a breakdown of the command
1. `--proto_path` is the directory of the `.proto` files
2. `--go_out=service` specifies the service directory as the root for the output
3. By default, the behavior is to place the output file in the directory named after the Go packages import path, such as one provided by the `go_package` option within the `.proto` file (can also be specified with the --go_opt=paths=import option). `--go_opt=module=productinfo/service`follows that logic, but removes the specified module from the output filename (`productinfo/service/ecommerce` becomes `ecommerce`).
4. The go-grpc flags serve the same purpose as the other flags, but for the grpc client and server output

The current project layout places the protobuf definitions under the chapter-x/proto directory. If you wished to keep the protobuf definitions in the service project itself (`service/proto`), then you would use this command from the `service` root: `protoc --proto_path=proto --go_out=. --go_opt=module=productinfo/service`.

Reuse these files for the client as well, placing them under the go/client/ecommerce folder.

### Proto 2023 message changes

Protobuf Edition 2023 introduces a more explicit and consistent way to handle optional fields. Before Edition 2023, unset optional fields would be represented by their zero value. This would cause issues when the zero value couldn't be determined to be optional by the message recipient. Edition 2023 represents values through pointers, so that their existance can be verified using a nil check. To set these pointer values, use the `google.golang.org/protobuf/proto` scalar constructors to create a new value and return a pointer.

## Java service

Gradle allows you to handle protoc as a project dependency, so you don't need to download protoc manually.

Gradle protobuf plugin will coordinate the gRPC and protobuf Java generation.

1. Specify the source directory for the protobuf files, as well as the generated Java files
```kotlin
sourceSets {
    main {
        java {
            srcDir("build/generated/source/proto/main/grpc")
            srcDir("build/generated/source/proto/main/java")
        }
        proto {
            srcDir(project.rootDir.resolve("../../proto"))
        }
    }
}
```
2. Configure protobuf plugin to use the protoc library and the gRPC gen Java plugin
```kotlin
protobuf {
    protoc {
        artifact = libs.protoc.compiler.get().toString()
    }
    plugins {
        id("grpc") {
            artifact = libs.protoc.gen.plugin.get().toString()
        }
    }
    generateProtoTasks {
        all().forEach { task ->
            task.plugins {
                id("grpc") {}
            }
        }
    }
}
```
3. Run `./gradlew build` to generate the gRPC and protobuf Java files. They will be output in `build/generated/source/proto/main/grpc` and `build/generated/source/proto/main/java`.