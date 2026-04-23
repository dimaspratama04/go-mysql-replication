package config

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Database struct {
	Mysql *gorm.DB
}

// func NewInitDatabase(primaryCfg DBConfig, replicaCfg DBConfig) (*Database, error) {
// 	// connect primary
// 	primaryDB, err := NewDB(primaryCfg)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect primary DB: %w", err)
// 	}

// 	// connect replica
// 	replicaDB, err := NewDB(replicaCfg)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to connect replica DB: %w", err)
// 	}

// 	return &Database{
// 		Primary:  primaryDB,
// 		Replicas: replicaDB,
// 	}, nil
// }

func NewInitDatabase(databaseConfig DBConfig) (*Database, error) {
	// connect primary
	dbs, err := NewDB(databaseConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect primary DB: %w", err)
	}

	// connect replica
	return &Database{
		Mysql: dbs,
	}, nil
}

func NewDB(cfg DBConfig) (*gorm.DB, error) {
	fmt.Printf("Connecting to DB at %s:%s...\n", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=UTC",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Logging handled by zerolog
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database, on %s:%s/%s: %w", cfg.Host, cfg.Port, cfg.Name, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	log.Info().
		Str("host", cfg.Host).
		Str("port", cfg.Port).
		Str("database", cfg.Name).
		Msg("Database connection established")

	return db, nil
}

func ProxyDBConfig() DBConfig {
	return DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "root"),
		Name:     getEnv("DB_NAME", "products_db"),
	}
}

func PrimaryDBConfig() DBConfig {
	return DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "root"),
		Name:     getEnv("DB_NAME", "products_db"),
	}
}

func ReplicaDBConfig() DBConfig {
	return DBConfig{
		Host:     getEnv("DB_REPLICA_HOST", "localhost"),
		Port:     getEnv("DB_REPLICA_PORT", "3307"),
		User:     getEnv("DB_REPLICA_USER", "root"),
		Password: getEnv("DB_REPLICA_PASSWORD", "root"),
		Name:     getEnv("DB_REPLICA_NAME", "products_db"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
