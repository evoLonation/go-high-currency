package service

type BuyItemsService struct {
	*context
}

func (p *BuyItemsService) EnterItems(OrderId int64, Barcode int64, Number int64) (result bool, err error) {
	result = true
	err = nil
	return
}
