package config

// HertzConfig 定义 Hertz 配置文件的结构体
type HertzConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

// MySQLConfig 定义 mysql 配置文件结构体
type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	DBname       string `mapstructure:"db_name"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

// VideoCRCRedisConfig 定义 VideoCRC 的redis 配置文件结构体
type VideoCRCRedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// VideoFRCRedisConfig 定义 VideoFRC 的redis 配置文件结构体
type VideoFRCRedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// UserInfoRCRedisConfig 定义 UserInfoRC 的redis 配置文件结构体
type UserInfoRCRedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// JWTConfig 定义 jwt 配置文件结构体
type JWTConfig struct {
	SigningKey  string `mapstructure:"signing_key"`
	IdentityKey string `mapstructure:"identity_key"`
}

// System 定义项目配置文件结构体
type System struct {
	HertzConfig           *HertzConfig           `mapstructure:"hertz"`
	MySQLConfig           *MySQLConfig           `mapstructure:"mysql"`
	VideoCRCRedisConfig   *VideoCRCRedisConfig   `mapstructure:"videoCRCRedis"`
	VideoFRCRedisConfig   *VideoFRCRedisConfig   `mapstructure:"videoFRCRedis"`
	UserInfoRCRedisConfig *UserInfoRCRedisConfig `mapstructure:"userInfoRCRedis"`
	JWTConfig             *JWTConfig             `mapstructure:"jwt"`
}
