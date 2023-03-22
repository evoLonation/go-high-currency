package service

import "github.com/go-redis/redis/v8"

type HighPriorityService struct {
	*context
	enterItemsRdb *redis.Client
}

func (p *HighPriorityService) EnterItems(orderId int64, barcode int64, number int64) (result bool, err error) {
	result = true
	err = nil
	return
}
