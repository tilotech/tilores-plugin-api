# TiloRes Plugin API

This project defines all the plugin APIs that can be implemented to customize
TiloRes.

## How does it work?

The plugin mechanism of TiloRes is based on
[Tilos go-plugin](https://github.com/tilotech/go-plugin) library.

The biggest advantages compared to the
[standard go plugin mechanism](https://pkg.go.dev/plugin) is, that the plugin
provider and plugin consumer can be compiled independently from each other (with
diverging go versions) - a main requirement, because the plugins may be
developed from the customers while the consumers are mostly developed by us.

Comparing to other plugin libraries, this one also works in serverless
environments.

## How to implement a plugin provider?

Each sub package provides a single interface that needs to be implemented. After
you created your implementation all you need to do is create an executable
out of it, typically done by implementing the `main` function in a similar way:

*example using the `dispatcher` interface*
```go
package main

import (
	"github.com/tilotech/go-plugin"
	"github.com/tilotech/tilores-plugin-api/dispatcher"
)

func main() {
	err := plugin.ListenAndServe(api.Provide(&MyDispatcherImpl))
	if err != nil {
		panic(err)
	}
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