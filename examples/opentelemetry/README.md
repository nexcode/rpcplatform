## Tracing systems preparation

In this example, we will connect two distributed tracing systems such as [Zipkin](https://zipkin.io) and [Jaeger](https://www.jaegertracing.io).
You can run some of these systems or even use something else. List of supported [exporters](https://pkg.go.dev/go.opentelemetry.io/otel/exporters).

For the test, let's run two docker containers:

```shell
docker run -d --name zipkin -p 9411:9411 openzipkin/zipkin
docker run -d --name jaeger -p 16686:16686 -p 4317:4317 jaegertracing/all-in-one
```

## Adding new settings

This example is very similar to [QuickStart](../quickstart), the main change is that new settings have been added:

```go
otlpExporter, err := otlptracegrpc.New(context.Background(),
    otlptracegrpc.WithEndpoint("localhost:4317"),
    otlptracegrpc.WithInsecure(),
)

if err != nil {
    panic(err)
}

zipkinExporter, err := zipkin.New("http://localhost:9411/api/v2/spans")
if err != nil {
    panic(err)
}

rpcp, err := rpcplatform.New("rpcplatform", etcdClient,
    options.OpenTelemetry("ServiceName", 1, otlpExporter, zipkinExporter),
)

if err != nil {
    panic(err)
}
```

## Launching this demo

```shell
cd examples/opentelemetry/server
go run .
```

```shell
cd examples/opentelemetry/client
go run .
```

## Let's see the tracing reports

- Open Zipkin UI in a browser: `http://localhost:9411`
- Open Jaeger UI in a browser: `http://localhost:16686`
