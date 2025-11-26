## Tracing systems preparation

In this example, we will connect two distributed tracing systems such as [Zipkin](https://zipkin.io) and [Jaeger](https://www.jaegertracing.io).
You can run some of these systems or even use something else. List of supported [exporters](https://pkg.go.dev/go.opentelemetry.io/otel/exporters).

For the test, let's run two docker containers with tracing systems:

```shell
docker run -d --name zipkin -p 9411:9411 openzipkin/zipkin
docker run -d --name jaeger -p 16686:16686 -p 4317:4317 jaegertracing/all-in-one
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

![Zipkin](zipkin.png)
![Jaeger](jaeger.png)
