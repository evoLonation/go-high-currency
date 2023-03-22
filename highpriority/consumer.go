package highpriority

import (
	"context"
	"encoding/json"
	"go-high-currency/config"
	"log"
	"sync"
	"time"

	"github.com/adjust/rmq/v5"
)

type ConsumerRoutine[REQ any, RES any] struct {
	*commonContext
}

func NewConsumerRoutine[REQ any, RES any](config *config.RedisConf) *ConsumerRoutine[REQ, RES] {
	return &ConsumerRoutine[REQ, RES]{
		commonContext: newContext(config),
	}
}

func (p *ConsumerRoutine[REQ, RES]) StartComsumer(handler func(req *REQ) (*RES, error)) {
	if err := p.queue.StartConsuming(1000, 100*time.Millisecond); err != nil {
		log.Fatal(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	p.queue.AddConsumerFunc("handler", func(delivery rmq.Delivery) {
		wg.Wait()
		req := &req[REQ]{}
		err := json.Unmarshal([]byte(delivery.Payload()), req)
		if err != nil {
			log.Fatal(err)
		}
		topic := req.Topic
		param := req.Param
		res := &res[RES]{}
		res.Result, err = handler(param)
		res.Err = err.Error()
		if err := p.rdb.Publish(context.Background(), topic, "done").Err(); err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
	})
}
