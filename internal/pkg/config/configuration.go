package config

import (
	stdlog "log"
	"path/filepath"

	bearConfig "github.com/nite-coder/blackbear/pkg/config"
	"github.com/nite-coder/blackbear/pkg/config/provider/env"
	"github.com/nite-coder/blackbear/pkg/config/provider/file"
	"github.com/nite-coder/blackbear/pkg/log"
	"github.com/nite-coder/blackbear/pkg/log/handler/console"
	"github.com/nite-coder/blackbear/pkg/log/handler/gelf"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var (
	// EnvPrefix 是環境變數的前墬
	EnvPrefix string
)

// LogSetting 用來設定 log 相關資訊
type LogSetting struct {
	Name             string
	Type             string
	MinLevel         string `mapstructure:"min_level"`
	ConnectionString string `mapstructure:"connection_string"`
}

// Database 用來提供連線的資料庫數據
type Database struct {
	Name             string
	ConnectionString string `mapstructure:"connection_string"`
	Type             string
	IsMigrated       bool `mapstructure:"is_migrated"`
}

// Configuration 用來代表 config 設定物件
type Configuration struct {
	rootDirPath string

	Env       string
	Mode      string
	Logs      []LogSetting
	Databases []Database
	Temporal  struct {
		Address string
	}
	Redis struct {
		ClusterMode     bool `mapstructure:"cluster_mode"`
		Addresses       []string
		Password        string
		MaxRetries      int `mapstructure:"max_retries"`
		PoolSizePerNode int `mapstructure:"pool_size_per_node"`
		DB              int
	}
	Jaeger struct {
		AdvertiseAddr string `mapstructure:"advertise_addr"`
	}
	Frontend struct {
		HTTPBind          string `mapstructure:"http_bind"`
		HTTPAdvertiseAddr string `mapstructure:"http_advertise_addr"`
	}
	Event struct {
		GRPCBind          string `mapstructure:"grpc_bind"`
		GRPCAdvertiseAddr string `mapstructure:"grpc_advertise_addr"`
	}
	Wallet struct {
		GRPCBind          string `mapstructure:"grpc_bind"`
		GRPCAdvertiseAddr string `mapstructure:"grpc_advertise_addr"`
	}
}

// New function 創建一個 configuration instance 出來
func New(fileName string) Configuration {
	cfg := Configuration{}

	bearConfig.RemoveAllPrividers()

	envProvider := env.New()
	bearConfig.AddProvider(envProvider)

	fileProvder := file.New()
	err := fileProvder.Load()
	if err != nil {
		panic(err)
	}

	bearConfig.AddProvider(fileProvder)
	err = bearConfig.Scan("", &cfg)

	if err != nil {
		stdlog.Fatal("config: yaml unmarshal error:", err)
	}

	return cfg
}

func (cfg Configuration) Path(path ...string) string {
	return filepath.Join(cfg.rootDirPath, filepath.Join(path...))
}

func (cfg Configuration) InitLogger(appID string) {
	// set up log target
	log.
		Str("app_id", appID).
		Str("env", cfg.Env).
		SaveToDefault()

	for _, target := range cfg.Logs {
		switch target.Type {
		case "console":
			clog := console.New()
			levels := log.GetLevelsFromMinLevel(target.MinLevel)
			log.AddHandler(clog, levels...)
		case "gelf":
			graylog := gelf.New(target.ConnectionString)
			levels := log.GetLevelsFromMinLevel(target.MinLevel)
			log.AddHandler(graylog, levels...)
		}
	}
}

// TracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func (cfg Configuration) TracerProvider(appID string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cfg.Jaeger.AdvertiseAddr)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(appID),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}
