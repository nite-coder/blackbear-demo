package config

import (
	"flag"
	"io/ioutil"
	stdlog "log"
	"os"
	"path/filepath"

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
	Name     string
	Username string
	Password string
	Address  string
	DBName   string
}

// Prometheus 用來設定 prometheus
type Prometheus struct {
	SubSystemName string `yaml:"sub_system_name"`
	Bind          string `yaml:"bind"`
	MetricsPath   string `yaml:"metrics_path"`
}

// Configuration 用來代表 config 設定物件
type Configuration struct {
	Env       string
	Mode      string
	Logs      []LogSetting
	Databases []Database
	Redis     struct {
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
	flag.Parse()
	cfg := Configuration{}

	rootDirPath := os.Getenv("STARTER_HOME")
	if rootDirPath == "" {
		//read and parse config file
		rootDirPathStr, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			stdlog.Fatalf("config: file error: %s", err.Error())
		}
		rootDirPath = rootDirPathStr
	}
	configPath := filepath.Join(rootDirPath, "configs", fileName)
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

	if EnvPrefix == "" {
		stdlog.Fatal("config: env prefix not set")
	}

	return cfg
}
