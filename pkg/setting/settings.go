package setting

import (
	"context"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type S3Config struct {
	Endpoint        string   `yaml:"endpoint" envconfig:"S3_ENDPOINT"`
	AccessKeyID     string   `yaml:"accessKeyID" envconfig:"S3_ACCESS_KEY_ID"`
	SecretAccessKey string   `yaml:"secretAccessKey" envconfig:"S3_SECRET_ACCESS_KEY"`
	Bucket          string   `yaml:"bucket" envconfig:"S3_BUCKET"`
	RegionName      string   `yaml:"regionName" envconfig:"S3_REGION_NAME"`
	UseSSL          bool     `yaml:"useSSL" envconfig:"S3_USE_SSL"`
	Buckets         []string `yaml:"buckets" envconfig:"S3_BUCKETS"`
}

type CorsConfig struct {
	AllowOrigins []string `yaml:"allowOrigins" envconfig:"CORS_ALLOW_ORIGINS"`
}

type Config struct {
	S3Config   S3Config   `yaml:"s3"`
	CorsConfig CorsConfig `yaml:"cors"`

	Context struct {
		Context context.Context    `yaml:"-" envconfig:"-"`
		Cancel  context.CancelFunc `yaml:"-" envconfig:"-"`
	} `yaml:"-"`
}

var Settings = &Config{}

func init() {
	SetupSettings()
}

// SetupSettings initialize the configuration instance
func SetupSettings() {

	err := envconfig.Process("", Settings)
	if err != nil {
		log.Fatalf("setting, fail to get from env': %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	Settings.Context.Context = ctx
	Settings.Context.Cancel = cancel
}
