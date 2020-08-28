package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	stdlog "log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/jasonsoft/log/v2"
	"github.com/jasonsoft/log/v2/handlers/console"
	"github.com/jasonsoft/log/v2/handlers/gelf"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/label"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"gopkg.in/yaml.v2"
	gormMySQL "gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	Name       string
	Username   string
	Password   string
	Address    string
	DBName     string
	Type       string
	IsMigrated bool `yaml:"is_migrated"`
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
	BFF struct {
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

func (cfg Configuration) InitDatabase(name string) (*gorm.DB, error) {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(180) * time.Second

	var connectionString string
	for _, database := range cfg.Databases {
		if strings.EqualFold(database.Name, name) {

			switch strings.ToLower(database.Type) {
			case "mysql":
				connectionString = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&multiStatements=true", database.Username, database.Password, database.Address, database.DBName)
			}

			// migrate database if needed
			if database.IsMigrated {
				path := cfg.Path("deployments", "database", database.DBName)
				path = filepath.ToSlash(path) // due to migrate package path issue on window os, therefore, we need to run this
				source := fmt.Sprintf("file://%s", path)
				migrateDBURL := fmt.Sprintf("%s://%s", database.Type, connectionString)

				m, err := migrate.New(
					source,
					"mysql://root:root@tcp(localhost:3306)/starter_db",
				)
				if err != nil {
					return nil, fmt.Errorf("db migration config is wrong. db_name: %s, source: %s, migrateDBURL: %s, error: %w", database.DBName, source, migrateDBURL, err)
				}

				err = m.Up()
				if err != nil && !errors.Is(err, migrate.ErrNoChange) {
					return nil, fmt.Errorf("db migration failed. db: %s, source: %s, migrateDBURL: %s, error: %w", database.DBName, source, migrateDBURL, err)
				}

				log.Infof("%s database was migrated", database.DBName)
			}
		}
	}

	var db *gorm.DB
	var err error
	err = backoff.Retry(func() error {
		db, err = gorm.Open(gormMySQL.Open(connectionString), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("main: database open failed: %w", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			return err
		}

		sqlDB.SetMaxIdleConns(150)
		sqlDB.SetMaxOpenConns(300)
		sqlDB.SetConnMaxLifetime(14400 * time.Second)

		err = sqlDB.Ping()
		if err != nil {
			return fmt.Errorf("main: database ping failed: %w", err)
		}

		return nil
	}, bo)

	if err != nil {
		return nil, fmt.Errorf("main: database connect err: %w", err)
	}

	return db, nil
}

// InitTracer creates a new trace provider instance and registers it as global trace provider.
func (cfg Configuration) InitTracer() func() {
	// Create and install Jaeger export pipeline
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint(cfg.Jaeger.AdvertiseAddr),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: "event",
			Tags: []label.KeyValue{
				label.String("version", "1.0"),
			},
		}),
		jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)
	if err != nil {
		log.Err(err).Fatal("install jaeger pipleline failed.")
	}

	return func() {
		flush()
	}
}
