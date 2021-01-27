package setting

import (
	"context"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type S3Config struct {
	Endpoint        string `yaml:"endpoint" envconfig:"S3_ENDPOINT"`
	AccessKeyID     string `yaml:"accessKeyID" envconfig:"S3_ACCESS_KEY_ID"`
	SecretAccessKey string `yaml:"secretAccessKey" envconfig:"S3_SECRET_ACCESS_KEY"`
	Bucket          string `yaml:"bucket" envconfig:"S3_BUCKET"`
}

type Config struct {
	S3Config S3Config `yaml:"s3"`

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
	var err error
	err = envconfig.Process("", Settings)
	if err != nil {
		log.Fatalf("setting, fail to get from env': %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	Settings.Context.Context = ctx
	Settings.Context.Cancel = cancel
}
