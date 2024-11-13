package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// var loadEnvOnce sync.Once

// func GetEnv(key string) string {
// 	loadEnvOnce.Do(func() {
// 		err := godotenv.Load()
// 		if err != nil {
// 			log.Printf("Ошибка загрузки .env: %v", err)
// 		}
// 	})
// 	return os.Getenv(key)
// }

var DefaultStorageConfig = StorageConfig{
	Host:     "localhost",
	Port:     "5432",
	Database: "postgres",
	Username: "postgres",
	Password: "zxcvB526",
}

type StorageConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

func NewClient(ctx context.Context, maxAttempts int, sc StorageConfig) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", sc.Username, sc.Password, sc.Host, sc.Port, sc.Database)
	fmt.Printf("Подключение к БД с DSN: %s\n", dsn)
	err = DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.Connect(ctx, dsn)
		if err != nil {
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		log.Printf("Ошибка при подключении: %v", err)

	}
	if pool == nil {
		return nil, fmt.Errorf("не удалось создать pool соединений")
	}

	return pool, nil
}

func DoWithTries(fn func() error, attemtps int, delay time.Duration) (err error) {
	for attemtps > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attemtps--

			continue
		}

		return nil
	}

	return
}
