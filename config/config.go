package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Conf struct {
	HttpServer   ServerConf       `yaml:"httpServer"`
	Service      ServiceConf      `yaml:"service"`
	HighPriority HighPriorityConf `yaml:"highPriority"`
}
type HighPriorityConf struct {
	EnterItems RedisConf `yaml:"enterItems"`
}
type ServiceConf struct {
	DataSource    string            `yaml:"dataSource"`
	ReplicationDB ReplicationDBConf `yaml:"replicationDB"`
	ShardingDB    ShardingDBConf    `yaml:"shardingDB"`
	RedisCluster  RedisClusterConf  `yaml:"redisCluster"`
	HttpServer    ServerConf        `yaml:"httpServer"`
}
type ReplicationDBConf struct {
	MasterSource string `yaml:"masterSource"`
	ReadSource   string `yaml:"readSource"`
}
type ShardingDBConf struct {
	DatabaseNumber int64    `yaml:"databaseNumber"`
	TableNumber    int64    `yaml:"tableNumber"`
	DataSources    []string `yaml:"dataSources"`
}
type RedisClusterConf struct {
	Redis      []RedisConf `yaml:"redis"`
	NodeNumber int64       `yaml:"nodeNumber"`
}
type RedisConf struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
}
type ServerConf struct {
	Port            string    `yaml:"port"`
	EnterItemsRedis RedisConf `yaml:"redisServer"`
}

var configFile string = "./etc/config.yaml"

func ParseConfig() (config Conf) {
	dirs, err := os.ReadDir("./etc")
	if err != nil {
		log.Fatal(errors.Wrap(err, "read directory error"))
	}
	var dirInfo string
	for _, dir := range dirs {
		dirInfo += dir.Name() + ", "
	}
	log.Printf("files: %s\n", dirInfo)

	log.Println("start read config file")
	content, err := os.ReadFile(configFile)
	// log.Print(string(content))
	if err != nil {
		log.Fatal(errors.Wrap(err, "read config file error"))
	}
	if err := yaml.Unmarshal(content, &config); err != nil {
		log.Fatal(errors.Wrap(err, "unmarshal config file error"))
	}
	conf, _ := json.Marshal(&config)
	log.Print(string(conf))
	return
}
