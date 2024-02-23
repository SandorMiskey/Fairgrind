// region: packages

package main

import (
	// standard
	"encoding/json"
	"fmt"
	"log/syslog"
	"models"

	// redirected
	"utils"
	// "models"

	// 3rd party
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/log"
	// "github.com/gofiber/fiber/v2/middleware/cors"
	// "github.com/gofiber/fiber/v2/middleware/logger"
	// "github.com/gofiber/fiber/v2/middleware/recover"
	// "github.com/gofiber/fiber/v2/middleware/requestid"
)

// endregion: packages
// region: globals

var (
	err error

	DB    *gorm.DB
	Env   map[string]string
	Cache *redis.Client

	// CacheMutex = sync.RWMutex{}
)

const (
	LOG_ERR    syslog.Priority = syslog.LOG_ERR
	LOG_NOTICE syslog.Priority = syslog.LOG_NOTICE
	LOG_INFO   syslog.Priority = syslog.LOG_INFO
	LOG_DEBUG  syslog.Priority = syslog.LOG_DEBUG
	LOG_EMERG  syslog.Priority = syslog.LOG_EMERG

	ENV_PREFIX         string = "CLR_"
	ENV_CACHE_PASSWORD string = ENV_PREFIX + "CACHE_PASSWORD"
	ENV_CACHE_HOST     string = ENV_PREFIX + "CACHE_HOST"
	ENV_DB_HOST        string = ENV_PREFIX + "DB_HOST"
	ENV_DB_PARAMS      string = ENV_PREFIX + "DB_PARAMS"
	ENV_DB_PASSWORD    string = ENV_PREFIX + "DB_PASSWORD"
	ENV_DB_PORT        string = ENV_PREFIX + "DB_PORT"
	ENV_DB_SCHEMA      string = ENV_PREFIX + "DB_SCHEMA"
	ENV_DB_USER        string = ENV_PREFIX + "DB_USER"
	ENV_LOG_SEVERITY   string = ENV_PREFIX + "LOG_SEVERITY"
	ENV_LOG_PREFIX     string = ENV_PREFIX + "LOG_PREFIX"
	ENV_MQ_EXCHANGE    string = ENV_PREFIX + "MQ_EXCHANGE"
	ENV_MQ_HOST        string = ENV_PREFIX + "MQ_HOST"
	ENV_MQ_PASSWORD    string = ENV_PREFIX + "MQ_PASSWORD"
	ENV_MQ_PORT        string = ENV_PREFIX + "MQ_PORT"
	ENV_MQ_ROUTING     string = ENV_PREFIX + "MQ_ROUTING"
	ENV_MQ_QUEUE       string = ENV_PREFIX + "MQ_QUEUE"
	ENV_MQ_USER        string = ENV_PREFIX + "MQ_USER"
)

// endregion

func main() {

	// region: read .env

	Env, err = godotenv.Read()
	utils.Panic(err)

	// endregion: read .env
	// region: logger

	utils.GetLogger(Env[ENV_LOG_SEVERITY], Env[ENV_LOG_PREFIX])
	defer utils.Logger.Close()
	logger := utils.Logger.Out

	// endregion
	// region: db

	// TODO: use "database/sql" for connection pool and fine tunning

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", Env[ENV_DB_USER], Env[ENV_DB_PASSWORD], Env[ENV_DB_HOST], Env[ENV_DB_PORT], Env[ENV_DB_SCHEMA], Env[ENV_DB_PARAMS])
	logger(LOG_DEBUG, dsn)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	utils.Panic(err)

	// err = DB.SetupJoinTable(&models.Project{}, "Grinders", &models.ProjectGrinder{})
	// 	panic(err)

	// endregion
	// region: redis

	Cache = redis.NewClient(&redis.Options{
		Addr:     Env[ENV_CACHE_HOST],
		Password: Env[ENV_CACHE_PASSWORD],
		DB:       0,
		Protocol: 3,
	})
	logger(LOG_DEBUG, Cache)

	// cacheCtx := context.Background()
	// err = Cache.Set(cacheCtx, "key", "value", 0).Err()
	// if err != nil {
	// 	panic(err)
	// }
	// val, err := Cache.Get(cacheCtx, "key").Result()
	// if err == redis.Nil {
	// 	fmt.Println("key2 does not exist")
	// } else if err != nil {
	// 	utils.Panic(err)
	// } else {
	// 	fmt.Println("key", val)
	// }

	// endregion: redis
	// region: mq

	mqUrl := "amqp://" + Env[ENV_MQ_USER] + ":" + Env[ENV_MQ_PASSWORD] + "@" + Env[ENV_MQ_HOST] + ":" + Env[ENV_MQ_PORT] + "/"
	logger(LOG_DEBUG, mqUrl)

	mqConn, err := amqp.Dial(mqUrl)
	utils.Panic(err, "failed to connect to mq")
	defer mqConn.Close()

	mqCh, err := mqConn.Channel()
	utils.Panic(err, "failed to open a channel")
	defer mqCh.Close()

	mqQueue, err := mqCh.QueueDeclare(
		Env[ENV_MQ_QUEUE], // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	utils.Panic(err, "failed to declare a queue")

	// don't dispatch a new message to a worker until it has processed and acknowledged the previous one
	err = mqCh.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	utils.Panic(err, "failed to set qoS")

	// bind the queue to the exchange
	err = mqCh.QueueBind(
		Env[ENV_MQ_QUEUE],    // name
		Env[ENV_MQ_ROUTING],  // routing key
		Env[ENV_MQ_EXCHANGE], // exchange
		false,                // noWait
		nil,                  // arguments
	)
	utils.Panic(err, "failed to bind queue to exchange")

	mqMsgs, err := mqCh.Consume(
		mqQueue.Name, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	utils.Panic(err, "failed to register a consumer")

	// endregion

	var forever chan struct{}

	go func() {
		for d := range mqMsgs {

			// region: new msg

			logger(LOG_DEBUG, "received mq message", string(d.Body))
			d.Ack(false)

			// endregion
			// region: parse msg

			var msg models.MqMsg

			err := json.Unmarshal(d.Body, &msg)
			if err != nil {
				logger(LOG_ERR, "error processing message into models.MqMsg", err)
				continue
			}
			logger(LOG_DEBUG, msg)

			// endregion
			// region: routes

			if msg.Database != Env[ENV_DB_SCHEMA] {
				logger(LOG_DEBUG, "msg.Database != Env[ENV_DB_SCHEMA]", msg.Database, Env[ENV_DB_SCHEMA])
				continue
			}

			switch msg.Table {
			case "clearing_batches":
				logger(LOG_DEBUG, "routing table match", msg.Table)

				// var rawData json.RawMessage
				// err = json.Unmarshal(d.Body, &rawData)
				// if err != nil {
				// 	logger(LOG_ERR, "error processing message data into json.RawMessage", err)
				// 	continue
				// }
				// logger(LOG_DEBUG, rawData)
			default:
				logger(LOG_DEBUG, "no routing match", msg.Table)
				continue
			}

			// endregion

		}
	}()

	logger(LOG_INFO, "waiting forever")
	<-forever

}
