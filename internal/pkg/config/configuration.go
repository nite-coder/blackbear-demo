package config

import (
	"io/ioutil"
	stdlog "log"
	"os"
	"path/filepath"

	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/log/v2/handlers/console"
	"github.com/jasonsoft/log/v2/handlers/gelf"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"gopkg.in/yaml.v2"
)

var (
	// EnvPrefix 是 GAM 環境變數的前墬
	EnvPrefix string
)

// LogSetting 用來設定 log 相關資訊
type LogSetting struct {
	Name             string `yaml:"name"`
	Type             string `yaml:"type"`
	MinLevel         string `yaml:"min_level"`
	ConnectionString string `yaml:"connection_string"`
}

// Database 用來提供連線的資料庫數據
type Database struct {
	Name             string
	ConnectionString string `yaml:"connection_string"`
	Type             string
	IsMigrated       bool `yaml:"is_migrated"`
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
		ClusterMode     bool     `yaml:"cluster_mode"`
		Addresses       []string `yaml:"addresses"`
		Password        string   `yaml:"password"`
		MaxRetries      int      `yaml:"max_retries"`
		PoolSizePerNode int      `yaml:"pool_size_per_node"`
		DB              int      `yaml:"db"`
	}
	Jaeger struct {
		AdvertiseAddr string `yaml:"advertise_addr"`
	}
	Frontend struct {
		HTTPBind          string `yaml:"http_bind"`
		HTTPAdvertiseAddr string `yaml:"http_advertise_addr"`
	}
	Event struct {
		GRPCBind          string `yaml:"grpc_bind"`
		GRPCAdvertiseAddr string `yaml:"grpc_advertise_addr"`
	}
	Wallet struct {
		GRPCBind          string `yaml:"grpc_bind"`
		GRPCAdvertiseAddr string `yaml:"grpc_advertise_addr"`
	}
}

// New function 創建一個 configuration instance 出來
func New(fileName string) Configuration {
	cfg := Configuration{}

	cfg.rootDirPath = os.Getenv("STARTER_HOME")
	if cfg.rootDirPath == "" {
		//read and parse config file
		rootDirPathStr, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			stdlog.Fatalf("config: file error: %s", err.Error())
		}
		cfg.rootDirPath = rootDirPathStr
	}

	//configPath := filepath.Join(rootDirPath, "configs", fileName)
	configPath := cfg.Path("configs", fileName)
	_, err := os.Stat(configPath)
	if err != nil {
		stdlog.Fatalf("config: file error: %s", err.Error())
	}

	// config exists
	file, err := ioutil.ReadFile(filepath.Clean(configPath))
	if err != nil {
		stdlog.Fatalf("config: read file error: %s", err.Error())
	}

	err = yaml.Unmarshal(file, &cfg)
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

// InitTracer creates a new trace provider instance and registers it as global trace provider.
func (cfg Configuration) InitTracer(appID string) func() {
	// Create and install Jaeger export pipeline
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint(cfg.Jaeger.AdvertiseAddr),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: appID,
			Tags: []label.KeyValue{
				label.String("version", "1.0"),
			},
		}),
		jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)
	if err != nil {
		log.Err(err).Fatal("install jaeger pipleline failed.")
	}
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return func() {
		flush()
	}
}
