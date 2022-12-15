# About Config

In order to make it easier to initialize the SDK, we decided to encapsulate the `Option` that uses the configuration file on the basis of the code framework layer.

In order to make the configuration file format more common and form a standard configuration file, we design the configuration from the `Business` function level.

According to the standard, we use `.` for hierarchical distinction, and `Camel Case` for word distinction.

it is divided into two parts configuration:
1. Common configuration
2. Dedicated configuration

Example: [agent.yml](../example/agent.yml) 

## Common configuration
It is a configuration common to the entire SDK. These public configurations should have no dots `.`

There are currently the following configurations.

| config      | description                           | example             |
|-------------|---------------------------------------|---------------------|
| serviceName | string, the name of your service      | zone.damoin.service |
| address     | string, the sdk api host port address | 127.0.0.1:9900      |

## Dedicated configuration

It is a dedicated configuration for a certain business in the SDK. We use `.` for hierarchical distinction, and `Camel Case` for word distinction.

There are currently the following configurations.

| config                            | description                                                                     | example                            |
|-----------------------------------|---------------------------------------------------------------------------------|------------------------------------|
| tracing.type                      | string, the type of tracing                                                     | log-tracing                        |
| tracing.enable                    | bool, the tracing switch                                                        | true                               |
| tracing.sample.rate               | float64, the tracing sample rate, minimum=0,maximum=1                           | 1                                  |
| tracing.shared.spans              | bool, set the client to request whether the Span Id of the server uses the same | true                               |
| tracing.id128bit                  | bool, set the span id use 128 bit                                               | false                              |
| reporter.output.server            | string, Data sending service configuration                                      | http://localhost:9411/api/v2/spans |
| reporter.output.server.tls.enable | bool, whether the sending service needs to use tls certificate                  | false                              |
| reporter.output.server.tls.key    | string, the tls key of the output server                                        |                                    |
| reporter.output.server.tls.cert   | string, the tls cert of the output server                                       |                                    |
| reporter.output.server.tls.caCert | string, the tls ca cert of the output server                                    |                                    |