package service

import (
	"fmt"
	"go-high-currency/common"
	"go-high-currency/config"

	redis "github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type context struct {
	singleDB         *sqlx.DB
	shardingDB       []*sqlx.DB
	masterDB         *sqlx.DB
	readDB           *sqlx.DB
	rdbs             []*redis.Client
	shardingDBNum    int64
	shardingTableNum int64
	redisClusterNum  int64
}

func (p *context) shardingTableNameInteger(tableName string, id int64) (*sqlx.DB, string) {
	tableId := id % p.shardingTableNum
	dbId := id / p.shardingTableNum % p.shardingDBNum
	return p.shardingDB[dbId], fmt.Sprintf("%s_%d", tableName, tableId)
}
func (p *context) shardingTableNameString(tableName string, id string) (*sqlx.DB, string) {
	return p.shardingTableNameInteger(tableName, common.Hash(id))
}

func newContext(conf *config.ServiceConf) *context {

	shardingDB := make([]*sqlx.DB, conf.ShardingDB.DatabaseNumber)
	for i, source := range conf.ShardingDB.DataSources {
		shardingDB[i] = common.NewMysqlDB(source)
	}
	rdbs := make([]*redis.Client, conf.RedisCluster.NodeNumber)
	for i, redisConf := range conf.RedisCluster.Redis {
		rdbs[i] = common.NewRedisClient(&redisConf)
	}
	return &context{
		singleDB:         common.NewMysqlDB(conf.DataSource),
		shardingDB:       shardingDB,
		masterDB:         common.NewMysqlDB(conf.ReplicationDB.MasterSource),
		readDB:           common.NewMysqlDB(conf.ReplicationDB.ReadSource),
		rdbs:             rdbs,
		shardingDBNum:    conf.ShardingDB.DatabaseNumber,
		shardingTableNum: conf.ShardingDB.TableNumber,
		redisClusterNum:  conf.RedisCluster.NodeNumber,
	}
}

func NewServices(conf *config.ServiceConf) (buyItemsService *BuyItemsService) {
	context := newContext(conf)
	buyItemsService = &BuyItemsService{
		context: context,
	}
	return
}

func NewPriorityService(serviceConf *config.ServiceConf, highPriConf *config.HighPriorityConf) *HighPriorityService {
	context := newContext(serviceConf)
	return &HighPriorityService{
		context:       context,
		enterItemsRdb: common.NewRedisClient(&highPriConf.EnterItems),
	}
}
