This repo contains the start of an AWS VPC component written in Go: https://github.com/justinvp/vpc. This is stipp a WIP. It doesn't do much at the moment (it's not yet a full port of the Node.js `@pulumi/awsx` VPC component).

Note: The plan is to create a boilerplate repo similar to https://github.com/pulumi/pulumi-provider-boilerplate (or just replace this repo).

## Required changes in `pulumi/pulumi`

This VPC example requires some changes in [`pulumi/pulumi`](https://github.com/pulumi/pulumi) that have not been merged to `master` yet. The changes address the following:

1. [[codegen/python] Panic when referencing some external types](https://github.com/pulumi/pulumi/issues/5819)
2. [[codegen/python] Fix bugs referencing external resources/types](https://github.com/pulumi/pulumi/issues/6232)
3. [Support defining remote components in Go](https://github.com/pulumi/pulumi/issues/5489)

The WIP changes are in the https://github.com/pulumi/pulumi/commits/justin/goauthoring branch. I will be opening PRs for these imminently. (1) and (2) are ready, aside from some cleanup. (3) is WIP, but enough is there to be able to write `Construct` in Go, which this VPC example demonstrates.

## Authoring

Authoring is currently a very manual process.

### Schema authoring

The schema for the component needs to be defined manually. See https://github.com/justinvp/vpc/blob/main/provider/cmd/pulumi-gen-vpc/main.go

Note: We are considering making it so that the schema can be generated automatically from Go struct definitions (i.e. using `pulumi` tags and comments). In the meantime, it must be defined manually.

### Component authoring

Write your component as you would normally in Go. See https://github.com/justinvp/vpc/blob/main/provider/cmd/pulumi-resource-vpc/vpc.go

### Provider implementation

A provider with an implementation of the `Construct` gRPC method is necessary for other languages to be able to use the component. Inside `Construct` you'll create the instance of the component and return any resulting state.

There is a WIP change in `pulumi/pulumi` in the `justin/goauthoring` branch that provides some helper code to write a provider that makes it possible to write `Construct` in Go. Basically in your `func main()` you can call `pulumi.ProviderMain` and pass it the `Construct` function to use. This handles all the gRPC boilerplate, and handles converting the raw gRPC values into the Pulumi Go SDK model (e.g. provides a `pulumi.Context` to use, converts the input values to `Outputs` that properly track dependencies, and provides the options as a `pulumi.ResourceOption`).

See https://github.com/justinvp/vpc/blob/main/provider/cmd/pulumi-resource-vpc/main.go

Note: This design hasn't been discussed/reviewed yet, so there will likely be tweaks and changes based on feedback. Any feedback would be much appreciated!

## Building

To build:

```
make build
```

This will generate the schema, build the provider, and generate/build each language SDK.

The provider plugin binary is placed in `./bin`. You'll want to add this directory to `PATH` for local development for the CLI to find the plugin.

To use this from TypeScript, run:

```
make install_nodejs_sdk
```

Which will make the Node.js package yarn-linkable.

Similar for .NET:

```
make install_dotnet_sdk
```

For Python, you can run `python -m pip install -e` with the full path to the `./sdk/python/bin` directory.
