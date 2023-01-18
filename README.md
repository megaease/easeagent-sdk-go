# easeagent-sdk-go

A lightweight & opening Go SDK for Cloud-Native and APM system

- [easeagent-sdk-go](#easeagent-sdk-go)
  - [Overview](#overview)
    - [Principles](#principles)
  - [Features](#features)
  - [QuickStart](#quickstart)
    - [Init Agent](#init-agent)
      - [1. Get SDK](#1-get-sdk)
      - [2. Import package](#2-import-package)
      - [3. New Agent](#3-new-agent)
    - [Wrapping HTTP Server](#wrapping-http-server)
  - [Documentation](#documentation)
  - [Example](#example)
  - [About MegaEase Cloud](#about-megaease-cloud)
  - [Community](#community)
  - [Licenses](#licenses)

## Overview

- EaseAgent SDK can collect distributed application tracing, which could be used in the APM system and improve the observability of a distributed system. for the tracing, EaseAgent SDK follows the [Google Dapper](https://research.google/pubs/pub36356/) paper and use [zipkin-go](https://github.com/openzipkin/zipkin-go) core library. 
- EaseAgent SDK also can work with Cloud-Native architecture. For example, it can help Service Mesh (especially for [EaseMesh](https://github.com/megaease/easemesh/) ) to do some control panel work.
- EaseAgent SDK also can work with [MegaEase Cloud](https://cloud.megaease.com/). For example, it can monitor for service by Go Docker APP.

### Principles
- Safe to Go application/service.
- Lightweight and very low CPU, memory, and I/O resource usage.
- Highly extensible, users can easily do extensions through the api
- Design for Micro-Service architecture, collecting the data from a service perspective.

## Features
* Easy to use. It is right out of the box for HTTP Server Tracing.
  * Collecting Tracing Logs.
    * HTTP Server
    * HTTP Client
    * Supplying the `health check` endpoint
  * Decorate the Span API for Middleware

* Data Reports
  * Console Reporter.
  * HTTP Reporter.

* Standardization
    * The tracing data format is fully compatible with the Zipkin data format.

## QuickStart
### Init Agent
#### 1. Get SDK
```shell
go get github.com/megaease/easeagent-sdk-go
```
#### 2. Import package

```go
import (
    "github.com/megaease/easeagent-sdk-go/agent"
    "github.com/megaease/easeagent-sdk-go/plugins/zipkin"
)
```
#### 3. New Agent
```go
// new tracing agent from yaml file and sets host and port of Span.localEndpoint
// By default, use yamlFile="" is use easemesh.DefaultSpec() and Console Reporter for tracing.
// By default, use localHostPort="" is not set host and port of Span.localEndpoint.
var easeagent, _ = agent.NewWithOptions(agent.WithYAML(os.Getenv("EASEAGENT_CONFIG"), ":8090"))
```
### Wrapping HTTP Server
```go
func main() {
	// initialize router
	router := http.NewServeMux()
	http.ListenAndServe(":8090", easeagent.WrapUserHandler(router))
}
```

## Documentation
[About Config](./doc/about-config.md)

## Example

1. [HTTP example](./example/http/main.go)

2. [mesh example](./example/mesh/main.go)

## About MegaEase Cloud 
1. [Use SDK in MegaEase Cloud](./doc/how-to-use.md)
2. Get MegaEase Cloud Config. [About MegaEase Cloud Config](./doc/megaease-cloud-config.md)
3. [Decorate the Span](./doc/middleware-span.md). please use api: `zipkin.Tracing.StartMWSpan` and `zipkin.Tracing.StartMWSpanFromCtx` for decorate Span.

## Community

* [Github Issues](https://github.com/megaease/easeagent-php-go/issues)
* [Join Slack Workspace](https://join.slack.com/t/openmegaease/shared_invite/zt-upo7v306-lYPHvVwKnvwlqR0Zl2vveA) for requirement, issue and development.
* [MegaEase on Twitter](https://twitter.com/megaease)

If you have any questions, welcome to discuss them in our community. Welcome to join!


## Licenses
EaseAgent Go SDK is licensed under the Apache License, Version 2.0. See [LICENSE](./LICENSE) for the full license text.

