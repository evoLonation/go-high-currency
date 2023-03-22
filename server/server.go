package listen

import (
	"go-high-currency/common"
	"go-high-currency/config"
	"go-high-currency/highpriority"
	"go-high-currency/service"

	"github.com/gin-gonic/gin"
)

type enterItemsReq struct {
	OrderId int64 `json:"order_id"`
	Barcode int64 `json:"barcode"`
	Number  int64 `json:"number"`
}

type enterItemRes struct {
	Result bool `json:"result"`
}

func Start(serverConf *config.ServerConf, highConf *config.HighPriorityConf, buyItemsService *service.BuyItemsService, highPriorityService *service.HighPriorityService) {
	router := gin.Default()
	router.POST("buy-items-service/enter-items", common.GetGinHandler(func(req *enterItemsReq) (*enterItemRes, error) {
		ret, err := buyItemsService.EnterItems(req.OrderId, req.Barcode, req.OrderId)
		return &enterItemRes{Result: ret}, err
	}))
	producer := highpriority.NewProducer[enterItemsReq, enterItemRes](&highConf.EnterItems)
	router.POST("buy-items-service/enter-items/high-priority", common.GetGinHandler(func(req *enterItemsReq) (*enterItemRes, error) {
		return producer.Publish(req)
	}))
	go func() {
		router.Run(":" + serverConf.Port)
	}()
}
