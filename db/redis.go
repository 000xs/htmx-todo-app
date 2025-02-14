package db

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func Connect() (*redis.Client, *context.Context) {
	opt, err := redis.ParseURL("rediss://default:Ad5MAAIjcDE4N2VlOTczNmRkZmQ0ZWM4OGM5ZTdhODMyZGY2ZWFiZXAxMA@alert-redfish-56908.upstash.io:6379")
	if err != nil {
		panic(err)

	}
	fmt.Printf("%vRedis connected!%v", "\033[31m", "\033[31m\n")
	client := redis.NewClient(opt)

	return client, &ctx

}
