package data

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/beltran/gohive"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"github.com/shencw/kratos-realworld/internal/pkg/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewRealWorldDB,
	NewRedisConn,
	NewHiveConn,
	NewKafkaClient,
	NewAuthRepo,
	NewAccountRepo,
)

// Data .
type Data struct {
	realWorldDB *gorm.DB
	redisCli    *redis.Client
	hiveConn    *gohive.Connection
	kafkaCli    sarama.Client
}

// NewData .
func NewData(
	c *conf.Data,
	logger log.Logger,
	realWorldDB *gorm.DB,
	redisCli *redis.Client,
	hiveConn *gohive.Connection,
	kafkaCli sarama.Client,
) (*Data, func(), error) {
	cleanup := func() {
		logHelper := log.NewHelper(logger)
		if err := redisCli.Close(); err != nil {
			logHelper.Errorf("redis client connect close error: %s", err.Error())
		}
		if err := hiveConn.Close(); err != nil {
			logHelper.Errorf("hive connect close error: %s", err.Error())
		}
		if err := kafkaCli.Close(); err != nil {
			logHelper.Errorf("kafka client connect close error: %s", err.Error())
		}
		logHelper.Debug("closing the data resources")
	}
	return &Data{realWorldDB: realWorldDB, redisCli: redisCli, hiveConn: hiveConn, kafkaCli: kafkaCli}, cleanup, nil
}

// NewRealWorldDB .
func NewRealWorldDB(confServer *conf.Server, confData *conf.Data, logger log.Logger) *gorm.DB {
	logHelper := log.NewHelper(log.With(logger, "module", "NewRealWorldDB"))
	db, err := gorm.Open(mysql.Open(confData.GetDatabase().GetRealworld()), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		logHelper.Fatalf("failed to connect database: %v", err)
	}
	// 仅测试使用
	if confServer.GetEnv() == "prod" {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&Users{}); err != nil {
			logHelper.Fatalf("database auto migrate error: %v", err)
		}
	}
	logHelper.Info("mysql connect success: real_world")
	return db
}

// NewRedisConn .
func NewRedisConn(conf *conf.Data, logger log.Logger) *redis.Client {
	logHelper := log.NewHelper(log.With(logger, "module", "NewRedisCmd"))
	client := redis.NewClient(&redis.Options{
		Addr:         conf.Redis.Addr,
		ReadTimeout:  conf.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: conf.Redis.WriteTimeout.AsDuration(),
		DialTimeout:  time.Second * 2,
		PoolSize:     10,
	})
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()
	if err := client.Ping(timeout).Err(); err != nil {
		logHelper.Fatalf("redis connect error: %v", err)
	}
	logHelper.Infof("Connect Redis Success Addr:%s", conf.Redis.Addr)
	return client
}

// NewHiveConn .
func NewHiveConn(conf *conf.Data, logger log.Logger) *gohive.Connection {
	logHelper := log.NewHelper(log.With(logger, "module", "NewHiveConn"))
	configuration := gohive.NewConnectConfiguration()
	configuration.Username = conf.GetHive().GetConfiguration().GetUsername()
	configuration.Password = conf.GetHive().GetConfiguration().GetPassword()
	configuration.Database = conf.GetHive().GetConfiguration().GetDatabase()
	hiveConnect, err := gohive.Connect(conf.GetHive().GetHost(), int(conf.GetHive().GetPort()), conf.GetHive().GetAuth(), configuration)
	if err != nil {
		logHelper.WithContext(context.Background()).Fatalf("Hive connection error:%s", err)
		return nil
	}
	logHelper.WithContext(context.Background()).Infof("Hive connection success:%s:%d", conf.GetHive().GetHost(), conf.GetHive().GetPort())
	return hiveConnect
}

// NewKafkaClient .
// link: https://www.jianshu.com/p/666d2604e8f8
func NewKafkaClient(conf *conf.Data, logger log.Logger) sarama.Client {
	logHelper := log.NewHelper(log.With(logger, "module", "NewKafkaProducer"))
	c := sarama.NewConfig()
	c.Producer.Return.Successes = true // 同步模式
	client, err := sarama.NewClient(conf.GetKafka().GetAddress(), c)
	if err != nil {
		logHelper.Fatalf("unable to create kafka client: %q", err)
	}
	return client
}
