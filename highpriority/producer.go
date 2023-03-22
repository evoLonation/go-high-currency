package highpriority

import (
	"context"
	"encoding/json"
	"fmt"
	"go-high-currency/config"
	"log"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type Producer[REQ any, RES any] struct {
	*commonContext
}

func NewProducer[REQ any, RES any](config *config.RedisConf) *Producer[REQ, RES] {
	return &Producer[REQ, RES]{
		commonContext: newContext(config),
	}
}

func (p *Producer[REQ, RES]) Publish(paramStruct *REQ) (*RES, error) {
	topic := strconv.Itoa(int(time.Now().Unix()))
	ctx := context.Background()
	sub := p.rdb.Subscribe(ctx, topic)
	req := &req[REQ]{
		Topic: topic,
		Param: paramStruct,
	}
	payload, err := json.Marshal(req)
	if err != nil {
		log.Fatal(errors.Wrap(err, "json marshal error"))
	}
	if err := p.queue.Publish(string(payload)); err != nil {
		log.Printf("failed to publish: %s", err)
	}
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		log.Fatal(err)
	}
	res := &res[RES]{}
	if err := json.Unmarshal([]byte(msg.Payload), res); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Received message from " + msg.Channel + " channel.")
	sub.Unsubscribe(ctx)
	return res.Result, errors.New(res.Err)
}
