## Tracing system preparation

In this example, we will connect two distributed tracing systems: [Zipkin](https://zipkin.io) and [Jaeger](https://www.jaegertracing.io).
You can run these systems or use other alternatives. See the list of supported [exporters](https://pkg.go.dev/go.opentelemetry.io/otel/exporters).

For testing, let's run two Docker containers with the tracing systems:

```shell
docker run -d --name zipkin -p 9411:9411 openzipkin/zipkin
docker run -d --name jaeger -p 16686:16686 -p 4317:4317 jaegertracing/all-in-one
```

## Launching the demo

```shell
cd examples/opentelemetry/server
go run .
```

```shell
cd examples/opentelemetry/client
go run .
```

## Viewing tracing reports

- Open Zipkin UI in a browser: `http://localhost:9411`
- Open Jaeger UI in a browser: `http://localhost:16686`

![Zipkin](zipkin.png)
![Jaeger](jaeger.png)
