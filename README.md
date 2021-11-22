# TiloRes Plugin API

This project defines all the plugin APIs that can be implemented to customize
TiloRes.

## How does it work?

The plugin mechanism of TiloRes is based on
[Hashicorps go-plugin](https://github.com/hashicorp/go-plugin) library.

The biggest advantages compared to the
[standard go plugin mechanism](https://pkg.go.dev/plugin) is, that the plugin
provider and plugin consumer can be compiled independently from eachother (with
diverging go versions) - a main requirement, because the plugins may be
developed from the customers while the consumers are mostly developed from us.

> **Note:**
>
> Currently the communication uses `net/rpc` and not the also supported `gRPC`
> implementation. However, we will probably change this before the Plugin API
> hits the `v1` release. This will allow the plugin providers to provide non-go
> based implementations.

## How to implement a plugin provider?

Each sub package provides a single interface that needs to be implemented. After
you have created your implementation all you need to do is create an executable
out of it, typically done by implementing the `main` function in a similar way:

*example using the `dispatcher` interface*
```go
package main

import (
	"github.com/hashicorp/go-plugin"
	"github.com/tilotech/tilores-plugin-api/dispatcher"
)

func main() {
	var pluginMap = map[string]plugin.Plugin{
		"dispatcher": &dispatcher.Plugin{
			Impl: &MyDispatcherImpl{},
		},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: dispatcher.Handshake,
		Plugins:         pluginMap,
	})
}
```

## Where is it used?

The following list is intended to give an overview about plugin providers and
plugin consumers for each interface.

### [dispatcher](https://pkg.go.dev/github.com/tilotech/tilores-plugin-api/dispatcher#Dispatcher):

Known Providers:
* TiloRes Core Dispatcher - proprietary default implementation
* [Fake Dispatcher](https://github.com/tilotech/tilores-plugin-fake-dispatcher) -
  a dispatcher that fakes a TiloRes Core implementation for local API testing

Known Consumers:
* Customer projects generated using [TiloRes CLI](https://github.com/tilotech/tilores-cli)