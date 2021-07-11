module github.com/nite-coder/blackbear-demo

go 1.14

require (
	github.com/99designs/gqlgen v0.13.0
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.2.0
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/nite-coder/blackbear v0.0.0-20210710135651-97a27fc0a4df
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	github.com/vektah/gqlparser/v2 v2.2.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.21.0
	go.opentelemetry.io/otel v1.0.0-RC1
	go.opentelemetry.io/otel/bridge/opentracing v0.21.0
	go.opentelemetry.io/otel/exporters/jaeger v1.0.0-RC1
	go.opentelemetry.io/otel/sdk v1.0.0-RC1
	go.opentelemetry.io/otel/trace v1.0.0-RC1
	go.temporal.io/api v1.4.1-0.20210420220407-6f00f7f98373
	go.temporal.io/sdk v1.8.0
	golang.org/x/net v0.0.0-20210420210106-798c2154c571
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/tools v0.1.4 // indirect
	google.golang.org/grpc v1.39.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.1.1
	gorm.io/gorm v1.21.11
)
