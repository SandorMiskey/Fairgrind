// region: packages {{{

package main

import (
	// standard
	"encoding/json"
	"fmt"
	"log/syslog"
	"sort"

	// redirected
	"models"
	"utils"

	// 3rd party
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	// "github.com/redis/go-redis/v9"
)

// endregion: packages }}}
// region: globals {{{

var (
	env     map[string]string        = utils.GetEnv()
	json_pp func(interface{}) string = utils.JsonPP
	// Cache *redis.Client
	// CacheMutex = sync.RWMutex{}
)

const (
	// DB_CLEARING_BATCHES         string = "clearing_batches"
	// DB_CLEARING_BATCH_STATUS_ID string = "clearing_batch_status_id"
	// DB_CLEARING_BATCH_TYPE_ID   string = "clearing_batch_type_id"
	// DB_CLEARING_TASKS           string = "clearing_tasks"
	// DB_CLEARING_BATCH_ID     string = "clearing_batch_id"
	// DB_CLEARING_CLEARING_TASK_ID string = "clearing_clearing_task_id"

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

	ERR_DB_FAILED_TO_FETCH_BATCH_TYPE   string = "failed to fetch batch type"
	ERR_DB_FAILED_TO_FETCH_BATCH_STATUS string = "failed to fetch batch status"
	ERR_MQ_FAILED_TO_BIND               string = "failed to bind queue to exchange"
	ERR_MQ_FAILED_TO_CONNECT            string = "failed to connect to mq"
	ERR_MQ_FAILED_TO_CONSUME            string = "failed to register a consumer"
	ERR_MQ_FAILED_TO_DECLARE            string = "failed to declare a queue"
	ERR_MQ_FAILED_TO_PARSE              string = "failed to parse message"
	ERR_MQ_FAILED_TO_OPEN               string = "failed to open a channel"
	ERR_MQ_FAILED_TO_SET_QOS            string = "failed to set qoS"

	LOG_ERR    syslog.Priority = syslog.LOG_ERR
	LOG_NOTICE syslog.Priority = syslog.LOG_NOTICE
	LOG_INFO   syslog.Priority = syslog.LOG_INFO
	LOG_DEBUG  syslog.Priority = syslog.LOG_DEBUG
	LOG_EMERG  syslog.Priority = syslog.LOG_EMERG

	MSG_MSG_RECEIVED                   string = "received mq message"
	MSG_MSG_PARSED                     string = "parsed message"
	MSG_NO_ROUTING_MATCH               string = "no routing match"
	MSG_ROUTING_MATCH                  string = "routing match"
	MSG_SCHEMA_MISMATCH                string = "schema mismatch msg.Database != env[ENV_DB_SCHEMA]"
	MSG_SUBROUTING_BATCH_NONCREDITABLE string = "sub-routing exception: batch is non-creditable"
	MSG_SUBROUTING_MULTIPIER_ZERO      string = "sub-routing exception: new batch type multiplier is zero"
	MSG_SUBROUTING_MULTIPLIERS_MATCH   string = "sub-routing exception: old and new multipliers are the same"
	MSG_WAITING_FOREVER                string = "waiting forever"
)

// endregion }}}

func main() { // {{{

	// region: logger {{{

	utils.GetLogger(env[ENV_LOG_SEVERITY], env[ENV_LOG_PREFIX])
	defer utils.Logger.Close()
	logger := utils.Logger.Out

	// endregion }}}
	// region: db {{{

	// TODO: use "database/sql" for connection pool and fine tunning
	// TODO: use "sqlx" for better query handling

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", env[ENV_DB_USER], env[ENV_DB_PASSWORD], env[ENV_DB_HOST], env[ENV_DB_PORT], env[ENV_DB_SCHEMA], env[ENV_DB_PARAMS])
	logger(LOG_DEBUG, dsn, env)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	utils.Panic(err)

	// endregion }}}
	// region: redis {{{

	/*
		Cache = redis.NewClient(&redis.Options{
			Addr:     env[ENV_CACHE_HOST],
			Password: env[ENV_CACHE_PASSWORD],
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
	*/

	// endregion: redis }}}
	// region: mq {{{

	mqUrl := "amqp://" + env[ENV_MQ_USER] + ":" + env[ENV_MQ_PASSWORD] + "@" + env[ENV_MQ_HOST] + ":" + env[ENV_MQ_PORT] + "/"
	logger(LOG_DEBUG, mqUrl)

	mqConn, err := amqp.Dial(mqUrl)
	utils.Panic(err, ERR_MQ_FAILED_TO_CONNECT)
	defer mqConn.Close()

	mqCh, err := mqConn.Channel()
	utils.Panic(err, ERR_MQ_FAILED_TO_OPEN)
	defer mqCh.Close()

	mqQueue, err := mqCh.QueueDeclare(
		env[ENV_MQ_QUEUE], // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	utils.Panic(err, ERR_MQ_FAILED_TO_DECLARE)

	// don't dispatch a new message to a worker until it has processed and acknowledged the previous one
	err = mqCh.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	utils.Panic(err, ERR_MQ_FAILED_TO_SET_QOS)

	// bind the queue to the exchange
	err = mqCh.QueueBind(
		env[ENV_MQ_QUEUE],    // name
		env[ENV_MQ_ROUTING],  // routing key
		env[ENV_MQ_EXCHANGE], // exchange
		false,                // noWait
		nil,                  // arguments
	)
	utils.Panic(err, ERR_MQ_FAILED_TO_BIND)

	mqMsgs, err := mqCh.Consume(
		mqQueue.Name, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	utils.Panic(err, ERR_MQ_FAILED_TO_CONSUME)

	// endregion: mq }}}
	// region: forever {{{

	var forever chan struct{}

	go func() {
		for d := range mqMsgs {

			// region: new msg {{{

			logger(LOG_DEBUG, MSG_MSG_RECEIVED, string(d.Body))
			d.Ack(false)

			// endregion }}}
			// region: parse msg {{{

			var msg models.MqMsg

			err := json.Unmarshal(d.Body, &msg)
			if err != nil {
				logger(LOG_ERR, ERR_MQ_FAILED_TO_PARSE, err)
				continue
			}
			logger(LOG_DEBUG, msg.Xid, MSG_MSG_PARSED, json_pp(msg))

			// endregion }}}
			// region: routes {{{

			if msg.Database != env[ENV_DB_SCHEMA] {
				logger(LOG_DEBUG, msg.Xid, MSG_SCHEMA_MISMATCH, msg.Database, env[ENV_DB_SCHEMA])
				continue
			}

			// Order msg.Old keys alphabetically. We could also consider
			// pre-loading the current state of corresponding records at runtime
			// using GORM's pre-load feature when various routing rules occur,
			// but we prefer to use the data received from CDC, as we want the
			// execution to be as deterministic as possible. All in all, this
			// GORM approach might not have been the best idea, and should be
			// considered for rewriting to sqlx when the opportunity arises.
			//
			// The basic assumption is that credits once credited cannot be
			// deleted, and credits already marked as withdrawable cannot

			fields := make([]string, 0, len(msg.Old))
			for k := range msg.Old {
				fields = append(fields, k)
			}
			sort.Strings(fields)

			for _, field := range fields {

				// TODO: clearing_batch_type_id update {{{
				//
				// _MESSAGE_
				// msg.Type: update
				// msg.Table: clearing_batches
				// field: clearing_batch_type_id
				//
				// _CONDITIONS_
				// value points to a record in clearing_batch_types where the
				// multiplier is not equal to zero and differs from the previous
				// state.
				//
				// _ACTION_
				// Taking into account the batch status and the multiplier:
				// - If there are clearing_ledger records associated with the
				//   clearing_tasks records linked to the batch, then update
				//   those records.
				// - Create new clearing_ledger records for those clearing_tasks
				//   records associated with the batch for which a clearing_ledger
				//   event has not yet been registered."

				if msg.Type == "update" && msg.Table == "clearing_batches" && field == "clearing_batch_type_id" {
					logger(LOG_DEBUG, msg.Xid, MSG_ROUTING_MATCH, msg.Type, msg.Table, field, msg.Old[field])

					// validate types {{{

					var result *gorm.DB
					var id_old = uint(msg.Old[field].(float64))
					var id_new = uint(msg.Data[field].(float64))
					var type_old = models.ClearingBatchType{GORM: models.GORM{ID: id_old}}
					var type_new = models.ClearingBatchType{GORM: models.GORM{ID: id_new}}

					result = db.First(&type_new)
					if result.Error != nil {
						logger(LOG_ERR, msg.Xid, ERR_DB_FAILED_TO_FETCH_BATCH_TYPE, result.Error.Error())
						continue
					}

					// new multiplier is zero
					if type_new.Multiplier == 0 {
						logger(LOG_DEBUG, msg.Xid, MSG_SUBROUTING_MULTIPIER_ZERO)
						continue
					}

					result = db.First(&type_old)
					if result.Error != nil {
						logger(LOG_ERR, msg.Xid, ERR_DB_FAILED_TO_FETCH_BATCH_TYPE, result.Error.Error())
						continue
					}

					// type changed but multiplier is the same
					if type_old.Multiplier == type_new.Multiplier {
						logger(LOG_DEBUG, msg.Xid, MSG_SUBROUTING_MULTIPLIERS_MATCH)
						continue
					}

					// }}}
					// batch status {{{

					var status_id = uint(msg.Data["clearing_batch_status_id"].(float64))
					var status models.ClearingBatchStatus = models.ClearingBatchStatus{Id: status_id}
					result = db.Preload("ClearingLedgerStatus").First(&status)
					if result.Error != nil {
						logger(LOG_ERR, msg.Xid, ERR_DB_FAILED_TO_FETCH_BATCH_STATUS, result.Error.Error())
						continue
					}
					logger(LOG_DEBUG, msg.Xid, json_pp(status))

					// non-creditable batch status
					if status.ClearingLedgerStatus.Id == 0 {
						logger(LOG_DEBUG, msg.Xid, MSG_SUBROUTING_BATCH_NONCREDITABLE)
						continue
					}

					// }}}
					// TODO: records already in the ledger {{{

					logger(LOG_DEBUG, msg.Xid, "dskfléajsdféksjadfkésajféasdjkfsaédifjsaéifjaséd")
					continue

					// }}}
					// TODO: tasks not yet recorded in the ledger {{{

					// }}}

				}

				// }}}
				// TODO: clearing_batch_status_id update {{{
				//
				// _MESSAGE_
				// msg.Type: update
				// msg.Table: clearing_batches
				// field: clearing_batch_status_id
				//
				// _CONDITIONS_
				// The clearing_ledger_status_id corresponding to the
				// clearing_batch_statuses.id indicated by the value is not
				// NULL.
				//
				// _ACTION_
				// The same as when there is a change in the
				// clearing_batches.clearing_batch_type_id.

				// }}}
				// TODO: new record in clearing_tasks {{{
				//
				// _MESSAGE_
				// msg.Type: insert
				// msg.Table: clearing_tasks
				//
				// _ACTION_
				// Taking into account the batch status and type, creating a
				// ledger record corresponding to the task status and the
				// grinder-specific fee.

				// }}}
				// TODO: clearing_tasks.clearing_task_status_id update {{{
				//
				// _MESSAGE_
				// msg.Type: update
				// msg.Table: clearing_tasks
				// field: clearing_task_status_id
				//
				// _ACTION_
				// Taking into account the batch status and type, creating a
				// ledger record corresponding to the task status and the
				// grinder-specific fee.

				// }}}

				logger(LOG_DEBUG, msg.Xid, MSG_NO_ROUTING_MATCH, msg.Table, field, msg.Data[field])

			}

			/*
				// No performance difference between switch and if-else, it is purely
				// for aesthetics and code readability. I prefer switch over if-else
				// because it is easier to read, but you might think otherwise...

				switch msg.Table {
				case DB_CLEARING_BATCHES:
					logger(LOG_DEBUG, msg.Xid, MSG_ROUTING_MATCH, msg.Table)

					for field, value := range msg.Old {
						switch field {
						case DB_CLEARING_BATCH_STATUS_ID:
							logger(LOG_DEBUG, msg.Xid, MSG_SUBROUTING_MATCH, msg.Table+"."+field, value, msg.Data[field])
						default:
							logger(LOG_DEBUG, msg.Xid, MSG_NO_SUBROUTING_MATCH, msg.Table+"."+field)
						}
					}
				default:
					logger(LOG_DEBUG, msg.Xid, MSG_NO_ROUTING_MATCH, msg.Table)
					continue
				}
			*/

			// endregion }}}

		}
	}()
	logger(LOG_INFO, MSG_WAITING_FOREVER)

	<-forever

	// endregion: forever }}}

} // }}
