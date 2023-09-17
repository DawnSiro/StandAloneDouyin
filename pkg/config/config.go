package config

// HertzConfig 定义 Hertz 配置文件的结构体
type HertzConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

// MySQLConfig 定义 mysql 配置文件结构体
type MySQLConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	DBName      string `mapstructure:"dbName"`
	MaxOpenConn int    `mapstructure:"maxOpenConn"`
	MaxIdleConn int    `mapstructure:"maxIdleConn"`
}

// VideoRCRedisConfig 定义 VideoRedisClient 的 redis 配置文件结构体
type VideoRCRedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"poolSize"`
}

// CommentRCRedisConfig 定义 VideoCRedisClient 的 redis 配置文件结构体
type CommentRCRedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"poolSize"`
}

// UserRCRedisConfig 定义 UserInfoRedisClient 的 redis 配置文件结构体
type UserRCRedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"poolSize"`
}

// MessageRCRedisConfig 定义 消息模块 RedisClient 的 redis 配置文件结构体
type MessageRCRedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"poolSize"`
}

// JWTConfig 定义 jwt 配置文件结构体
type JWTConfig struct {
	SigningKey  string `mapstructure:"signingKey"`
	IdentityKey string `mapstructure:"identityKey"`
}

// PulsarConfig 定义 pulsar 配置文件结构体
type PulsarConfig struct {
	URL               string `mapstructure:"url"`
	ConnectionTimeout int64  `mapstructure:"connectTimeout"`
	OperationTimeout  int64  `mapstructure:"operatorTimeout"`
}

// System 定义项目配置文件结构体
type System struct {
	HertzConfig          *HertzConfig          `mapstructure:"hertz"`
	MySQLConfig          *MySQLConfig          `mapstructure:"mysql"`
	VideoRCRedisConfig   *VideoRCRedisConfig   `mapstructure:"videoRCRedis"`
	CommentRCRedisConfig *CommentRCRedisConfig `mapstructure:"commentRCRedis"`
	UserRCRedisConfig    *UserRCRedisConfig    `mapstructure:"userRCRedis"`
	MessageRCRedisConfig *MessageRCRedisConfig `mapstructure:"messageRCRedis"`
	JWTConfig            *JWTConfig            `mapstructure:"jwt"`
	PulsarConfig         *PulsarConfig         `mapstructure:"pulsar"`
}
