package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Web           *WebConfig           `mapstructure:"web"`
	Chat          *ChatConfig          `mapstructure:"chat"`
	Uploader      *UploaderConfig      `mapstructure:"uploader"`
	Observability *ObservabilityConfig `mapstructure:"observability"`
}

type WebConfig struct {
	Http struct {
		Port string
	}
}

type ChatConfig struct {
	Http struct {
		Port    string
		MaxConn int64
	}
	Redis struct {
		Password       string
		Addrs          string
		ExpirationHour int64
	}
	Message struct {
		MaxNum      int64
		MaxSizeByte int64
		Worker      int
	}
	JWT struct {
		Secret           string
		ExpirationSecond int64
	}
	Match struct {
		Worker int
	}
}

type UploaderConfig struct {
	Http struct {
		Port string
	}
	S3 struct {
		Endpoint   string
		DisableSSL bool
		Region     string
		Bucket     string
		AccessKey  string
		SecretKey  string
	}
	JWT struct {
		Secret string
	}
}

type ObservabilityConfig struct {
	Prometheus struct {
		Port string
	}
	Tracing struct {
		JaegerUrl string
	}
}

func setDefault() {
	viper.SetDefault("web.http.port", "5000")

	viper.SetDefault("chat.http.port", "5001")
	viper.SetDefault("chat.http.maxConn", 200)
	viper.SetDefault("chat.redis.password", "")
	viper.SetDefault("chat.redis.addrs", "localhost:6379")
	viper.SetDefault("chat.redis.expirationHour", 24)
	viper.SetDefault("chat.message.maxNum", 500)
	viper.SetDefault("chat.message.maxSizeByte", 4096)
	viper.SetDefault("chat.message.worker", 4)
	viper.SetDefault("chat.jwt.secret", "replaceme")
	viper.SetDefault("chat.jwt.expirationSecond", 86400)
	viper.SetDefault("chat.match.worker", 4)

	viper.SetDefault("uploader.http.port", "5002")
	viper.SetDefault("uploader.s3.endpoint", "http://localhost:9000")
	viper.SetDefault("uploader.s3.disableSSL", false)
	viper.SetDefault("uploader.s3.region", "us-east-1")
	viper.SetDefault("uploader.s3.bucket", "myfilebucket")
	viper.SetDefault("uploader.s3.accessKey", "")
	viper.SetDefault("uploader.s3.secretKey", "")
	viper.SetDefault("uploader.jwt.secret", "replaceme")

	viper.SetDefault("observability.prometheus.port", "8080")
	viper.SetDefault("observability.tracing.jaegerUrl", "")
}

func NewConfig() (*Config, error) {
	setDefault()

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
