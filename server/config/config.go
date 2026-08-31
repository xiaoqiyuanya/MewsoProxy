package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	App      AppConfig      `mapstructure:"app"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	User        string `mapstructure:"user"`
	Password    string `mapstructure:"password"`
	DBName      string `mapstructure:"db_name"`
	Charset     string `mapstructure:"charset"`
	AutoMigrate bool   `mapstructure:"auto_migrate"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	AccessSecret  string `mapstructure:"access_secret"`
	RefreshSecret string `mapstructure:"refresh_secret"`
	AccessTTL     string `mapstructure:"access_ttl"`
	RefreshTTL    string `mapstructure:"refresh_ttl"`
}

type AppConfig struct {
	RegisterEnabled bool   `mapstructure:"register_enabled"`
	DefaultGroupID  int    `mapstructure:"default_group_id"`
	SubscribeURL    string `mapstructure:"subscribe_url"`
	SSHEncryptKey   string `mapstructure:"ssh_encrypt_key"`
}

func Load() *Config {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetEnvPrefix("MEWSO")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	v.SetDefault("server.port", 8080)
	v.SetDefault("database.host", "127.0.0.1")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.user", "root")
	v.SetDefault("database.password", "")
	v.SetDefault("database.charset", "utf8mb4")
	v.SetDefault("database.auto_migrate", false)
	v.SetDefault("redis.addr", "127.0.0.1:6379")
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)
	v.SetDefault("jwt.access_secret", "change-me-access")
	v.SetDefault("jwt.refresh_secret", "change-me-refresh")
	v.SetDefault("jwt.access_ttl", "30m")
	v.SetDefault("jwt.refresh_ttl", "168h")
	v.SetDefault("app.register_enabled", true)
	v.SetDefault("app.default_group_id", 1)
	v.SetDefault("app.subscribe_url", "http://localhost:8081")
	v.SetDefault("app.ssh_encrypt_key", "")

	if _, err := os.Stat("config.yaml"); err == nil {
		_ = v.ReadInConfig()
	}

	if err := v.Unmarshal(&Config{}); err != nil {
		panic(fmt.Sprintf("config unpack failed: %v", err))
	}
	return build(v)
}

func build(v *viper.Viper) *Config {
	return &Config{
		Server:   ServerConfig{Port: v.GetInt("server.port")},
		Database: dbCfg(v),
		Redis:    RedisConfig{Addr: v.GetString("redis.addr"), Password: v.GetString("redis.password"), DB: v.GetInt("redis.db")},
		JWT:      JWTConfig{AccessSecret: v.GetString("jwt.access_secret"), RefreshSecret: v.GetString("jwt.refresh_secret"), AccessTTL: v.GetString("jwt.access_ttl"), RefreshTTL: v.GetString("jwt.refresh_ttl")},
		App:      AppConfig{RegisterEnabled: v.GetBool("app.register_enabled"), DefaultGroupID: v.GetInt("app.default_group_id"), SubscribeURL: v.GetString("app.subscribe_url"), SSHEncryptKey: v.GetString("app.ssh_encrypt_key")},
	}
}

func dbCfg(v *viper.Viper) DatabaseConfig {
	return DatabaseConfig{
		Host:        v.GetString("database.host"),
		Port:        v.GetInt("database.port"),
		User:        v.GetString("database.user"),
		Password:    v.GetString("database.password"),
		DBName:      v.GetString("database.db_name"),
		Charset:     v.GetString("database.charset"),
		AutoMigrate: v.GetBool("database.auto_migrate"),
	}
}

func (c *Config) Validate() error {
	if c.Server.Port <= 0 {
		return fmt.Errorf("server.port 缺失，请检查 config.yaml 或环境变量 MEWSO_SERVER_PORT")
	}
	if c.Database.Host == "" || c.Database.DBName == "" {
		return fmt.Errorf("database.host / database.db_name 缺失，请检查 config.yaml")
	}
	if c.Redis.Addr == "" {
		return fmt.Errorf("redis.addr 缺失，请检查 config.yaml")
	}
	return nil
}
