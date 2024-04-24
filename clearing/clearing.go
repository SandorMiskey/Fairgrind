// packages {{{

package main

import (
	// standard
	"encoding/json"
	"fmt"
	"sync"
	"time"

	// redirected
	"models"
	"utils"

	// 3rd party
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// packages }}}
// globals {{{

var (
	Env    map[string]string             = utils.GetEnv()
	Db     *gorm.DB                      = nil
	JsonPP func(interface{}) string      = utils.JsonPP
	Lock   sync.Mutex                    = sync.Mutex{}
	Logger func(...interface{}) *[]error = nil

	// Cache *redis.Client
	// CacheMutex = sync.RWMutex{}
)

// }}}

func main() { // {{{

	// logger {{{

	utils.GetLogger(Env[ENV_LOG_SEVERITY], Env[ENV_LOG_PREFIX])
	defer utils.Logger.Close()
	Logger = utils.Logger.Out

	// endregion }}}
	// db {{{

	// TODO: use "database/sql" for connection pool and fine tunning
	// TODO: use "sqlx" for better query handling. All in all, this GORM
	// approach might not have been the best idea, and should be considered for
	// rewriting to sqlx when the opportunity arises.

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", Env[ENV_DB_USER], Env[ENV_DB_PASSWORD], Env[ENV_DB_HOST], Env[ENV_DB_PORT], Env[ENV_DB_SCHEMA], Env[ENV_DB_PARAMS])
	Logger(LOG_DEBUG, dsn, Env)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	utils.Panic(err)
	Db = db
	Logger(LOG_DEBUG, Db)

	// endregion }}}
	// redis {{{

	/*
		Cache = redis.NewClient(&redis.Options{
			Addr:     Env[ENV_CACHE_HOST],
			Password: Env[ENV_CACHE_PASSWORD],
			DB:       0,
			Protocol: 3,
		})
		Logger(LOG_DEBUG, Cache)

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

	// redis }}}
	// mq {{{

	mqUrl := "amqp://" + Env[ENV_MQ_USER] + ":" + Env[ENV_MQ_PASSWORD] + "@" + Env[ENV_MQ_HOST] + ":" + Env[ENV_MQ_PORT] + "/"
	Logger(LOG_DEBUG, mqUrl)

	mqConn, err := amqp.Dial(mqUrl)
	utils.Panic(err, ERR_MQ_FAILED_TO_CONNECT)
	defer mqConn.Close()

	mqCh, err := mqConn.Channel()
	utils.Panic(err, ERR_MQ_FAILED_TO_OPEN)
	defer mqCh.Close()

	mqQueue, err := mqCh.QueueDeclare(
		Env[ENV_MQ_QUEUE], // name
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
		Env[ENV_MQ_QUEUE],    // name
		Env[ENV_MQ_ROUTING],  // routing key
		Env[ENV_MQ_EXCHANGE], // exchange
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

	// mq }}}
	// forever {{{

	var forever chan struct{}

	// on message received {{{

	go func() {
		defer func() {
			if r := recover(); r != nil {
				Logger(LOG_ERR, r)
			}
			Lock.Unlock()
		}()

		for d := range mqMsgs {

			// parse msg {{{

			var msg models.MqMsg

			d.Ack(false)
			Logger(LOG_INFO, MSG_MSG_RECEIVED)

			err := json.Unmarshal(d.Body, &msg)
			if err != nil {
				Logger(LOG_ERR, ERR_MQ_FAILED_TO_PARSE, string(d.Body), err)
				continue
			}
			Logger(LOG_DEBUG, msg.Xid, MSG_MSG_PARSED, utils.JsonPP(msg))

			// endregion }}}
			// schema mismatch {{{

			if msg.Database != Env[ENV_DB_SCHEMA] {
				Logger(LOG_DEBUG, msg.Xid, MSG_SCHEMA_MISMATCH, msg.Database, Env[ENV_DB_SCHEMA])
				continue
			}

			// }}}
			// memo {{{

			// The basic assumption is that credits once credited cannot be
			// deleted, and credits already marked as withdrawable cannot be
			// reclassified during the clearing process.
			//
			// BTW, no performance difference between switch and if-else, it is
			// purely for aesthetics and code readability. I prefer switch over
			// if-else because it is easier to read, but you might think
			// otherwise...

			// }}}
			// routes {{{

			switch {

			// new record in clearing_tasks {{{
			//
			// _MESSAGE_
			// msg.Type: insert
			// msg.Table: clearing_tasks
			//
			// _ACTION_
			// Taking into account the batch status and type, creating a
			// ledger record corresponding to the task status and the
			// grinder-specific fee.

			case msg.Type == "insert" && msg.Table == "clearing_tasks":
				Logger(LOG_INFO, msg.Xid, MSG_ROUTING_MATCH, msg.Type, msg.Table)
				ClearTask(uint(msg.Data["id"].(float64)))

			// }}}
			// default {{{
			default:
				Logger(LOG_INFO, msg.Xid, MSG_ROUTING_MISMATCH, msg.Type, msg.Table)

				// }}}

			}
			/* obsolete - to be deleted {{{
			switch msg.Type {

			// updates {{{

			case "update":

				// order msg.Old keys alphabetically {{{

				// Order msg.Old keys alphabetically. It is to be considered later
				// whether this is appropriate in this way, considering the cases
				// where, for example, the type and status of a batch are changed
				// in a single UPDATE. We could also consider pre-loading the
				// current state of corresponding and related records at runtime
				// using GORM's pre-load feature when various routing rules occur,
				// but we prefer to use the data received from CDC, as we want the
				// execution to be as deterministic as possible.

				fields := []string{}
				for k := range msg.Old {
					if k == "updated_at" || k == "created_at" || k == "deleted_at" {
						Logger(LOG_DEBUG, msg.Xid, MSG_SUBROUTING_SKIPPING_FIELD, k)
						continue
					}
					fields = append(fields, k)
				}
				sort.Strings(fields)

				// }}}

				for _, field := range fields {

					// TODO: clearing_batches.clearing_batch_type_id update {{{
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

					if msg.Table == "clearing_batches" && field == "clearing_batch_type_id" {
						Logger(LOG_INFO, msg.Xid, MSG_ROUTING_MATCH, msg.Type, msg.Table, field)

						// batch type {{{

						var result *gorm.Db
						var id_old = uint(msg.Old[field].(float64))
						var id_new = uint(msg.Data[field].(float64))
						var type_old = models.ClearingBatchType{GORM: models.GORM{ID: id_old}}
						var type_new = models.ClearingBatchType{GORM: models.GORM{ID: id_new}}

						result = db.First(&type_old)
						if result.Error != nil {
							Logger(LOG_ERR, msg.Xid, ERR_DB_FAILED_TO_FETCH_BATCH_TYPE, result.Error.Error())
							continue
						}

						result = db.First(&type_new)
						if result.Error != nil {
							Logger(LOG_ERR, msg.Xid, ERR_DB_FAILED_TO_FETCH_BATCH_TYPE, result.Error.Error())
							continue
						}

						// sub-routing exception: new multiplier is either less than or equal to the old one
						if type_old.Multiplier >= type_new.Multiplier {
							Logger(LOG_INFO, msg.Xid, MSG_SUBROUTING_MULTIPLIERS_MISMATCH, type_old.Multiplier, type_new.Multiplier)
							continue
						}

						// }}}
						// batch status {{{

						var status_id = uint(msg.Data["clearing_batch_status_id"].(float64))
						var status = models.ClearingBatchStatus{Id: status_id}

						result = db.Preload("ClearingLedgerStatus").First(&status)
						if result.Error != nil {
							Logger(LOG_ERR, msg.Xid, ERR_DB_FAILED_TO_FETCH_BATCH_STATUS, result.Error.Error())
							continue
						}
						Logger(LOG_INFO, msg.Xid, utils.JsonPP(status))

						// non-creditable batch status
						// if status.ClearingLedgerStatus.Id == 0 {
						// 	Logger(LOG_DEBUG, msg.Xid, MSG_SUBROUTING_BATCH_NONCREDITABLE)
						// 	continue
						// }

						// }}}
						// TODO: records already in the ledger {{{

						// }}}

						Logger(LOG_INFO, msg.Xid, MSG_SUBROUTING_DONE)
					}

					// }}}
					// TODO: clearing_batches.clearing_batch_status_id update {{{
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

				}

			// }}}

			}
			*/ // }}}

			Logger(LOG_INFO, msg.Xid, MSG_ROUTING_DONE)

			// }}}

		}
	}()

	// }}}
	// ticker {{{

	interval, err := time.ParseDuration(Env[ENV_TICKER_INTERVAL])
	utils.Panic(err)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				Logger(LOG_ERR, r)
			}
			Lock.Unlock()
		}()

		for range ticker.C {
			Logger(LOG_INFO, Env[ENV_TICKER_MARKER])
			Lock.Lock()
			// uncleared tasks {{{

			tasks := []models.ClearingTask{}
			qt := Db.Find(&tasks, "cleared_at IS NULL")
			if qt.Error != nil {
				Logger(LOG_ERR, qt.Error)
			} else {
				Logger(LOG_INFO, MSG_TASK_UNCLEARED, len(tasks))
				for _, task := range tasks {
					ClearTask(task.ID)
				}
			}

			// }}}
			// TODO: clearing_batches & clearing_ledger sync
			Lock.Unlock()
		}
	}()

	// }}}

	Logger(LOG_INFO, MSG_WAITING_FOREVER)

	<-forever

	// forever }}}

} // }
