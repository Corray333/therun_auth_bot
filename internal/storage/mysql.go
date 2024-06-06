package storage

import (
	"fmt"
	"os"

	redis "github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	DB    *sqlx.DB
	Redis *redis.Client
}

func NewStorage() *Storage {

	db := sqlx.MustConnect("mysql", os.Getenv("MYSQL_CONN_STR"))
	if err := db.Ping(); err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	if res := redisClient.Ping(); res.Err() != nil {
		panic(res.Err())
	}

	return &Storage{
		DB:    db,
		Redis: redisClient,
	}
}

func (store *Storage) SaveCode(phone, code string, type_id int) error {
	fmt.Println()
	fmt.Println(phone, type_id)
	fmt.Println()
	_, err := store.DB.Queryx(`
		INSERT INTO number_verification_requests (phone, type_id, request_id, is_active, code)
		VALUES (?, ?, NULL, 1, ?)
		ON DUPLICATE KEY UPDATE type_id    = VALUES(type_id),
								request_id = VALUES(request_id),
								is_active  = 1,
								code       = VALUES(code)
	`, phone, type_id, code)

	return err
}
