package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Web           *WebConfig           `mapstructure:"web"`
	Chat          *ChatConfig          `mapstructure:"chat"`
	Match         *MatchConfig         `mapstructure:"match"`
	Uploader      *UploaderConfig      `mapstructure:"uploader"`
	User          *UserConfig          `mapstructure:"user"`
	Kafka         *KafkaConfig         `mapstructure:"kafka"`
	Cassandra     *CassandraConfig     `mapstructure:"cassandra"`
	Redis         *RedisConfig         `mapstructure:"redis"`
	Observability *ObservabilityConfig `mapstructure:"observability"`
}

type WebConfig struct {
	Http struct {
		Server struct {
			Port string
		}
	}
}

type ChatConfig struct {
	Http struct {
		Server struct {
			Port    string
			MaxConn int64
			Swag    bool
		}
	}
	Grpc struct {
		Server struct {
			Port string
		}
	}
	Message struct {
		MaxNum        int64
		PaginationNum int
		MaxSizeByte   int64
	}
	JWT struct {
		Secret string
	}
}

type MatchConfig struct {
	Http struct {
		Server struct {
			Port    string
			MaxConn int64
			Swag    bool
		}
	}
	Grpc struct {
		Client struct {
			Chat struct {
				Endpoint string
			}
			User struct {
				Endpoint string
			}
		}
	}
	JWT struct {
		Secret           string
		ExpirationSecond int64
	}
}

type UploaderConfig struct {
	Http struct {
		Server struct {
			Port          string
			Swag          bool
			MaxBodyByte   int64
			MaxMemoryByte int64
		}
	}
	S3 struct {
		Endpoint  string
		Region    string
		Bucket    string
		AccessKey string
		SecretKey string
	}
}

type UserConfig struct {
	Http struct {
		Server struct {
			Port string
			Swag bool
		}
	}
	Grpc struct {
		Server struct {
			Port string
		}
	}
}

type KafkaConfig struct {
	Addrs string
}

type CassandraConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Keyspace string
}

type RedisConfig struct {
	Password                string
	Addrs                   string
	ExpirationHour          int64
	MinIdleConn             int
	PoolSize                int
	ReadTimeoutMilliSecond  int64
	WriteTimeoutMilliSecond int64
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
	viper.SetDefault("web.http.server.port", "5000")

	viper.SetDefault("chat.http.server.port", "5001")
	viper.SetDefault("chat.http.server.maxConn", 200)
	viper.SetDefault("chat.http.server.swag", false)
	viper.SetDefault("chat.grpc.server.port", "4000")
	viper.SetDefault("chat.message.maxNum", 5000)
	viper.SetDefault("chat.message.paginationNum", 5000)
	viper.SetDefault("chat.message.maxSizeByte", 4096)
	viper.SetDefault("chat.jwt.secret", "replaceme")

	viper.SetDefault("match.http.server.port", "5002")
	viper.SetDefault("match.http.server.maxConn", 200)
	viper.SetDefault("match.http.server.swag", false)
	viper.SetDefault("match.grpc.client.chat.endpoint", "localhost:4000")
	viper.SetDefault("match.grpc.client.user.endpoint", "localhost:4001")
	viper.SetDefault("match.jwt.secret", "replaceme")
	viper.SetDefault("match.jwt.expirationSecond", 86400)

	viper.SetDefault("uploader.http.server.port", "5003")
	viper.SetDefault("uploader.http.server.swag", false)
	viper.SetDefault("uploader.http.server.maxBodyByte", "67108864")   // 64MB
	viper.SetDefault("uploader.http.server.maxMemoryByte", "16777216") // 16MB
	viper.SetDefault("uploader.s3.endpoint", "http://localhost:9000")
	viper.SetDefault("uploader.s3.region", "us-east-1")
	viper.SetDefault("uploader.s3.bucket", "myfilebucket")
	viper.SetDefault("uploader.s3.accessKey", "")
	viper.SetDefault("uploader.s3.secretKey", "")

	viper.SetDefault("user.http.server.port", "5004")
	viper.SetDefault("user.http.server.swag", false)
	viper.SetDefault("user.grpc.server.port", "4001")

	viper.SetDefault("kafka.addrs", "localhost:9092")

	viper.SetDefault("cassandra.host", "localhost")
	viper.SetDefault("cassandra.port", 9042)
	viper.SetDefault("cassandra.user", "cassandra")
	viper.SetDefault("cassandra.password", "cassandra")
	viper.SetDefault("cassandra.keyspace", "randomchat")

	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.addrs", "localhost:6379")
	viper.SetDefault("redis.expirationHour", 24)
	viper.SetDefault("redis.minIdleConn", 30)
	viper.SetDefault("redis.poolSize", 500)
	viper.SetDefault("redis.readTimeoutMilliSecond", 500)
	viper.SetDefault("redis.writeTimeoutMilliSecond", 500)

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
