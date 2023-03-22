package main

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"strconv"
	"time"

	"github.com/pkg/errors"
	redis "github.com/redis/go-redis/v9"
)

type Order struct {
	GenerateId int64
	UserId     int64
}
type Sale struct {
	GenerateId int64
	Number     int64
	ItemId     int64
	OrderId    int64
}
type User struct {
	UserId   int64
	Username string
}
type Item struct {
	Id          int64
	StockNumber int64
	Price       float64
}

type ItemNeed struct {
	StockNumber int64 `json:"stock_number"`
}

func NewRedisClient(number int) (*redis.Client, context.Context) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("redis-cluster-%d.redis-cluster:6379", number),
		Password: "2002116yy", // 没有密码，默认值
		DB:       0,           // 默认DB 0
	})
	log.Println("get rdb success")
	ctx := context.Background()
	return rdb, ctx
}
func CacheToRedis(rdb *redis.Client, ctx context.Context) {
	item := &Item{
		Id:          1,
		StockNumber: 100,
		Price:       66.6,
	}
	itemValue, err := json.Marshal(&ItemNeed{
		StockNumber: item.StockNumber,
	})
	if err != nil {
		log.Fatal(errors.Wrap(err, "marshal itemValue error"))
	}
	itemKey := fmt.Sprintf("item-%d", item.Id)

	rdb.Set(ctx, itemKey, string(itemValue), time.Hour)

	log.Printf("append item : %s", string(itemValue))
}

func Handler(rdb *redis.Client, ctx context.Context, userId int64, itemId int64) (bool, error) {
	// definition
	var user *User = &User{}
	user.UserId = userId
	var itemValue *ItemNeed = &ItemNeed{}
	var itemValueStr string
	itemValueStr, err := rdb.Get(ctx, fmt.Sprintf("item-%d", itemId)).Result()
	// log.Printf("get item: %s\n", itemValueStr)
	if err != nil {
		log.Fatal(errors.Wrap(err, "item does not exists"))
	}
	if err := json.Unmarshal([]byte(itemValueStr), itemValue); err != nil {
		log.Fatal(errors.Wrap(err, "item unmarshal error"))
	}
	// precondition: 可以根据userId和itemId找到user和item
	// postcondition
	var sale *Sale = &Sale{}
	var order *Order = &Order{}
	if itemValue.StockNumber >= 1 {
		// 对已有对象的修改
		itemValue.StockNumber = itemValue.StockNumber - 1
		// 新对象的创建
		sale.GenerateId = time.Now().Unix()
		sale.ItemId = itemId
		sale.Number = 1
		order.GenerateId = time.Now().Unix()
		sale.OrderId = order.GenerateId
		order.UserId = userId
		// 插入order和sale，更新item
		itemValueByte, err := json.Marshal(itemValue)
		if err != nil {
			log.Fatal(errors.Wrap(err, "marshal error"))
		}
		if _, err := rdb.Set(ctx, fmt.Sprintf("item-%d", itemId), string(itemValueByte), time.Hour).Result(); err != nil {
			log.Fatal("append error")
		}
		return true, nil
	}
	return false, nil
}

func main() {
	rdb0, ctx0 := NewRedisClient(0)
	rdb1, ctx1 := NewRedisClient(1)
	rdb2, ctx2 := NewRedisClient(2)

	rdbs := []*redis.Client{rdb0, rdb1, rdb2}
	ctxs := []context.Context{ctx0, ctx1, ctx2}

	for i := 0; i < 10; i++ {
		key := "item-" + strconv.Itoa(i)
		id := redisSelector(key)
		log.Printf("set id : %d\n", id)
		if _, err := rdbs[id].Set(ctxs[id], key, string("i am "+key), time.Hour).Result(); err != nil {
			log.Fatal("append error")
		}
	}

	for i := 0; i < 10; i++ {
		key := "item-" + strconv.Itoa(i)
		id := redisSelector(key)
		log.Printf("get id : %d, value is\n", id)
		res, err := rdbs[id].Get(ctxs[id], key).Result()
		if err != nil {
			log.Println(errors.Wrap(err, "item does not exists"))
		}
		log.Println(res)
	}

	// for i := 0; i < 10001; i++ {
	// 	ret, err := Handler(rdb, ctx, 1, 1)
	// 	if err != nil {
	// 		log.Fatal("handler error")
	// 	}
	// 	if i%1000 == 0 {
	// 		log.Print("handler result : ")
	// 		println(ret)
	// 	}
	// }

}

func redisSelector(key string) int {
	h := fnv.New32()
	h.Write([]byte(key))
	return int(h.Sum32()) % 3
}

// package main

// import (
//     "crypto/md5"
//     "fmt"
//     "hash/fnv"
// )

// func main() {
//     strs := []string{"apple", "banana", "cherry", "durian"}
//     n := 3
//     h := fnv.New32()

//     for i := 0; i < n; i++ {
//         sum := uint32(0)
//         for _, s := range strs {
//             h.Write([]byte(s))
//             sum += h.Sum32()
//             h.Reset()
//         }
//         fmt.Printf("%d ", sum%10) // 将结果对10取模，得到0-9的数字
//     }
// }
// 在上面的代码中，我们选择了FNV哈希算法，将每个字符串转换为字节数组，并使用哈希
