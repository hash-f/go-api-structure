- [ ] 1. Add OpenTelemetry dependencies

  - [ ] 1.1. Add Go module dependencies: `go.opentelemetry.io/otel`, `go.opentelemetry.io/otel/sdk`, `go.opentelemetry.io/otel/exporters/otlp/otlptrace`, `go.opentelemetry.io/otel/exporters/otlp/otlpmetric`, `go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp`, `go.opentelemetry.io/contrib/instrumentation/database/sql/otelsql`
  - [ ] 1.2. Run `go mod tidy` and commit updated `go.mod` / `go.sum`

- [ ] 2. Create telemetry bootstrap package

  - [ ] 2.1. Create `internal/telemetry` package with `init.go` to set up providers
  - [ ] 2.2. Configure resource attributes (`service.name`, `service.version`, `deployment.environment`)
  - [ ] 2.3. Configure OTLP exporter with environment-configurable endpoint
  - [ ] 2.4. Expose `Shutdown(ctx)` to flush providers on graceful shutdown

- [ ] 3. Instrument HTTP server

  - [ ] 3.1. Wrap router/handlers with `otelhttp` middleware
  - [ ] 3.2. Ensure context propagation through request lifecycle
  - [ ] 3.3. Verify HTTP spans are exported to collector

- [ ] 4. Instrument SQL database

  - [ ] 4.1. Wrap database driver with `otelsql` and register as `otel-driver`
  - [ ] 4.2. Update `SQLStore` initialization to use instrumented driver
  - [ ] 4.3. Confirm query spans include operation name and latency

- [ ] 5. Add custom spans and metrics

  - [ ] 5.1. Identify critical business operations (e.g., user creation, transaction processing)
  - [ ] 5.2. Add spans around these operations with meaningful attributes
  - [ ] 5.3. Record custom metrics via `meter` (request count, DB latency)

- [ ] 6. Configure logging correlation

  - [ ] 6.1. Inject trace and span IDs into structured logs (zap)
  - [ ] 6.2. Verify logs correlate with traces in backend UI

- [ ] 7. Provide local OpenTelemetry Collector setup

  - [ ] 7.1. Add `docker-compose.yml` with `otel-collector`, Jaeger UI, Prometheus
  - [ ] 7.2. Include `otel-config.yaml` with OTLP receiver and exporters
  - [ ] 7.3. Document how to run collector locally

- [ ] 8. Update CI/CD pipeline

  - [ ] 8.1. Disable telemetry export during unit tests (use noop exporter)
  - [ ] 8.2. Add integration test to validate OTLP exporter connectivity
  - [ ] 8.3. Ensure pipeline spins up collector for end-to-end tests

- [ ] 9. Update documentation

  - [ ] 9.1. Add README section for telemetry setup and environment variables
  - [ ] 9.2. List required env vars (`OTEL_SERVICE_NAME`, `OTEL_EXPORTER_OTLP_ENDPOINT`, etc.)
  - [ ] 9.3. Provide troubleshooting tips for common issues

- [ ] 10. Review and cleanup
  - [ ] 10.1. Run `golangci-lint` to ensure lint passes
  - [ ] 10.2. Benchmark to assess performance impact of instrumentation
  - [ ] 10.3. Conduct code review and merge changes into main branch
