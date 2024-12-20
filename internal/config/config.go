package config

import (
	"casion/internal/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type AppConfig struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	JWTSecret     string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	ServerPort    string
}

var (
	DB     *gorm.DB
	Config AppConfig
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	Config = AppConfig{
		DBHost:        getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:        getEnvOrDefault("DB_PORT", "3306"),
		DBUser:        getEnvOrDefault("DB_USER", "root"),
		DBPassword:    getEnvOrDefault("DB_PASS", "Admin1234"),
		DBName:        getEnvOrDefault("DB_NAME", "casion_db"),
		JWTSecret:     getEnvOrDefault("JWT_SECRET", "ABCD1234567890ABCDabcd"),
		RedisHost:     getEnvOrDefault("REDIS_HOST", "localhost"),
		RedisPort:     getEnvOrDefault("REDIS_PORT", "6379"),
		RedisPassword: getEnvOrDefault("REDIS_PASSWORD", ""),
		ServerPort:    getEnvOrDefault("SERVER_PORT", "8080"),
	}

	log.Printf("DB Config - Host: %s, Port: %s, User: %s, Pass: %s, Name: %s",
		Config.DBHost, Config.DBPort, Config.DBUser, Config.DBPassword, Config.DBName)

	// Create DSN
	dsn := Config.DBUser + ":" + Config.DBPassword + "@tcp(" + Config.DBHost + ":" + Config.DBPort + ")/" + Config.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"

	// Open database connection
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migrate the schema
	if err := DB.AutoMigrate(&models.User{}, &models.Transaction{}); err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}

	log.Println("Database connected and migrated successfully")
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
