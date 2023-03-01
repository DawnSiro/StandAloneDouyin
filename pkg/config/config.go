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

// VideoCRCRedisConfig 定义 VideoCRC 的redis 配置文件结构体
type VideoCRCRedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"poolSize"`
}

// VideoFRCRedisConfig 定义 VideoFRC 的redis 配置文件结构体
type VideoFRCRedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"poolSize"`
}

// UserInfoRCRedisConfig 定义 UserInfoRC 的redis 配置文件结构体
type UserInfoRCRedisConfig struct {
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

// System 定义项目配置文件结构体
type System struct {
	HertzConfig           *HertzConfig           `mapstructure:"hertz"`
	MySQLConfig           *MySQLConfig           `mapstructure:"mysql"`
	VideoCRCRedisConfig   *VideoCRCRedisConfig   `mapstructure:"videoCRCRedis"`
	VideoFRCRedisConfig   *VideoFRCRedisConfig   `mapstructure:"videoFRCRedis"`
	UserInfoRCRedisConfig *UserInfoRCRedisConfig `mapstructure:"userInfoRCRedis"`
	JWTConfig             *JWTConfig             `mapstructure:"jwt"`
}
