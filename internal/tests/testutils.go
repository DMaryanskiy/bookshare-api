package tests

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/DMaryanskiy/bookshare-api/internal/db"
	"github.com/DMaryanskiy/bookshare-api/internal/db/models"
	"github.com/DMaryanskiy/bookshare-api/pkg/utils"
	"github.com/golang-migrate/migrate/v4"
	migrate_postgres "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func init() {
	envPath := "../../.env"
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Failed to load env: %v", err)
	}
}

var TestDB *gorm.DB

func SetupTestDB(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	dbConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	db.DB = dbConn
	TestDB = dbConn

	RunMigrations(t, dbConn)
}

func RunMigrations(t *testing.T, db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get raw DB: %v", err)
	}

	driver, err := migrate_postgres.WithInstance(sqlDB, &migrate_postgres.Config{})
	if err != nil {
		t.Fatalf("failed to create migrate driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations", // path to your migrations folder
		"postgres",
		driver,
	)
	if err != nil {
		t.Fatalf("failed to create migrate instance: %v", err)
	}

	// Firstly, downgrade DB
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("migration downgrade failed: %v", err)
	}

	// Now upgrading DB
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("migration upgrade failed: %v", err)
	}
}

var TestRedis *redis.Client

func SetupTestRedis() {
	addr := os.Getenv("TEST_REDIS_URL")
	TestRedis = redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   1,
	})

	// Optional: flush DB before test
	if err := TestRedis.FlushDB(context.Background()).Err(); err != nil {
		panic("failed to flush test Redis: " + err.Error())
	}
}

func GetAuthHeader(token string) string {
	return "Bearer " + token
}

func CreateTestUser(t *testing.T, email, password string) (models.User, string) {
	hashedPassword, _ := utils.HashPassword(password)
	user := models.User{
		Email:        email,
		PasswordHash: hashedPassword,
		IsVerified:   true,
	}
	require.NotNil(t, db.DB, "db.DB is nil — make sure SetupTestDB was called before this")
	err := db.DB.Table("auth.users").Create(&user).Error
	require.NoError(t, err)

	token, err := utils.GenerateAccessToken(user.ID.String(), 24*time.Hour)
	require.NoError(t, err)

	return user, token
}
