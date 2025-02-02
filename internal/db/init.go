package db

import (
	"context"
	"database/sql"
	"log"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

var (
	Mdb *sql.DB
	Rdb *redis.Client
	ctx = context.Background()
)

func InitMySql(dsn string) {
	var err error
	Mdb, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if err = Mdb.Ping(); err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}

	log.Println("Successfully connected to MySQL")
}

func InitRedis(RedisAddr string) {
	Rdb = redis.NewClient(&redis.Options{
		Addr: RedisAddr,
	})
	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis: %v", err)
	}
	log.Println("Successfully connected to Redis")
}
