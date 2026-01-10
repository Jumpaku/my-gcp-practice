package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
)

func main() {
	var ctx = context.Background()

	// 1. クライアントの初期化
	// ポート6379で動作しているValkeyに接続します
	var rdb redis.UniversalClient
	redisHost := os.Getenv("REDIS_HOST")
	switch {
	case redisHost == "":
		panic("REDIS_HOST environment variable is not set.")
	case strings.HasPrefix(redisHost, "localhost"):
		rdb = redis.NewClient(&redis.Options{
			Addr:     redisHost,
			PoolSize: 10,
		})
	default:
		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    []string{redisHost},
			PoolSize: 10,
		})
	}

	// 2. PINGコマンド (接続確認)
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("接続エラー: %v", err)
	}
	fmt.Printf("接続成功: %s\n", pong)

	// 3. SETコマンド (書き込み互換性確認)
	key := "test_key"
	value := "Hello Valkey from go-redis!"
	err = rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		log.Fatalf("書き込みエラー: %v", err)
	}
	fmt.Println("書き込み成功")

	// 4. GETコマンド (読み込み互換性確認)
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		log.Fatalf("読み込みエラー: %v", err)
	}
	fmt.Printf("読み込み成功: %s -> %s\n", key, val)

	// 5. INFOコマンド (サーバー情報の確認)
	// ここでサーバーが実際にValkeyであることを確認します
	info, err := rdb.Info(ctx, "server").Result()
	if err != nil {
		log.Fatalf("Info取得エラー: %v", err)
	}

	// INFOの結果から "redis_version" や "valkey_version" を探す
	fmt.Println("\n--- サーバー情報 ---")
	fmt.Println(info)
}
