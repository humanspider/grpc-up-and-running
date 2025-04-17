# Set Up

The current project includes the generated Protobuf and gRPC service+stub code, but you can regenerate them if changes are made to the project(s).

## Install protoc
`protoc` is the Protobuf compiler. This project requires generated code for both the Protobuf messages, and gRPC service+stub.

Download the latest release zip from here: https://github.com/protocolbuffers/protobuf/releases/

Uncompress the file `protobuf-<version>-<arch>.<tar.gz or zip>`.

Place the extracted contents of `bin` in `usr/local/bin` and `include` in `usr/local/include`.

## Golang generation

**NOTE:** Golang gRPC generation plugin currently supports editions, so we will use the `product_info_edition.proto` file.

In your service project directory, create a new Go module: `go mod init productinfo/service`

### Install gRPC library

Run `go get -u google.golang.org/grpc`

### Install protoc Go plugin

Run `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`and `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc` to install the plugin in `$GOBIN`(defaults to $GOPATH/bin). Include `$GOBIN` in your path so that protoc can find the plugins.

### Generate Go output

Navigate to the the root of the Golang service or client project. Run the compilation command: `protoc --proto_path=../../proto/edition --go_out=.--go_opt=module=productinfo/service product_info.proto --go-grpc_out=. --go-grpc_opt=module=productinfo/service`. This will output the generated Go protobuf inside of service/ecommerce.

Here's a breakdown of the command
1. `--proto_path` is the directory of the `.proto` files.
2. `--go_out=.` specifies current directory as the output for the message and base service definitions.
3. By default, the behavior is to place the output file in the directory named after the Go packages import path, such as one provided by the `go_package` option within the `.proto` file (can also be specified with the `--go_opt=paths=import option`). `--go_opt=module=productinfo/service`follows that logic, but removes the specified module from the output filename (`productinfo/service/ecommerce` becomes `ecommerce`).
4. The `go-grpc` flags serve the same purpose as the other flags, but for the gRPC stub and server output.

Reuse these files for the client as well, placing them under the go/client/ecommerce folder.

### Proto 2023 message changes

Protobuf Edition 2023 introduces a more explicit and consistent way to handle optional fields. Before Edition 2023, unset optional fields would be represented by their zero value. This would cause issues when the zero value couldn't be determined to be optional by the message recipient. Edition 2023 represents values through pointers, so that their existance can be verified using a nil check. To set these pointer values, use the `google.golang.org/protobuf/proto` scalar constructors to create a new value and return a pointer.

## Java Generation

**NOTE:** Java gRPC generation plugin currently does not support editions, so we will use the `product_info.proto` file.

### Option 1: gRPC generation plugin
Download the gRPC plugin from Maven repository https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/<version>/protoc-gen-grpc-java-<version>-<os-arch>.exe.
Example:https://repo1.maven.org/maven2/io/grpc/protoc-gen-grpc-java/1.71.0/protoc-gen-grpc-java-1.71.0-windows-x86_64.exe.

Rename it to `protoc-gen-grpc-java` and ensure that it has permissions to be executed.

Add the plugin binary to your path (recommend placing in `/usr/local/bin` for permanent use) or use the --plugin (recommend placing in a custom directory).

From the client/service project root, run `protoc --proto_path=../../proto product_info.proto --java_out=build/generated/source/proto/main/java --grpc-java_out=build/generated/source/proto/main/grpc`.

### Option 2: Gradle Protobuf code generation
Gradle allows you to handle protoc as a project dependency, so you don't need to download protoc manually.

Gradle protobuf plugin will coordinate the gRPC and protobuf Java generation.

1. Specify the source directory for the protobuf files, as well as tell the project where the generated Java files will be located.
```kotlin
sourceSets {
    main {
        java {
            srcDir("build/generated/source/proto/main/grpc")
            srcDir("build/generated/source/proto/main/java")
        }
        proto {
            srcDir(project.rootDir.resolve("../../proto"))
            exclude("*_edition.proto")
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
        all().configureEach { task ->
            task.plugins {
                id("grpc") {}
            }
        }
    }
}
```
This configuration adds the gRPC Java gen plugin and executes the tasks. This includes the `java` task builtin by default, and the explicit `grpc` task plugin.

Additional documentation can be found here: https://github.com/google/protobuf-gradle-plugin.

3. Run `./gradlew build` to generate the gRPC and protobuf Java files. They will be output in `build/generated/source/proto/main/grpc` and `build/generated/source/proto/main/java`.

### Gradle uber jar compilation

The `jar` task is included by default as part of the `java` plugin. By default, this will not include all of the external dependencies needed to run the program as a standalone jar, or uber jar. You must specify the extra files to include in the jar.

```kotlin
tasks.named<Jar>("jar") {
    manifest {
        attributes["Main-Class"] = "ecommerce.ProductInfoClient"
    }
    from(configurations.runtimeClasspath.get().map {
        if (it.isDirectory) it else zipTree(it)
    }) {
        duplicatesStrategy = DuplicatesStrategy.EXCLUDE
    }
}
```
* `from(..)` specifies the source files and directories to be included.
* `configurations.runtimeClasspath.get()` collects all of the runtime dependencies into a `FileCollection`.
* `if (it.isDirectory) it else zipTree(it)` leaves directories as-is and unpacks the JARs into individual classes and resources.