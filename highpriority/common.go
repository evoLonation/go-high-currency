package highpriority

import (
	"go-high-currency/common"
	"go-high-currency/config"
	"log"

	"github.com/adjust/rmq/v5"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type req[REQ any] struct {
	Topic string `json:"topic"`
	Param *REQ   `json:"param"`
}
type res[RES any] struct {
	Err    string `json:"err"`
	Result *RES   `json:"result"`
}

type commonContext struct {
	queue rmq.Queue
	rdb   *redis.Client
}

func newContext(config *config.RedisConf) *commonContext {
	rdb := common.NewRedisClient(config)
	connection, err := rmq.OpenConnectionWithRedisClient("consumer", rdb, nil)
	if err != nil {
		log.Fatal(errors.Wrap(err, "open rmq connection error"))
	}
	queue, err := connection.OpenQueue("queue")
	if err != nil {
		log.Fatal(errors.Wrap(err, "open rmq queue error"))
	}
	return &commonContext{
		queue: queue,
		rdb:   rdb,
	}
}
