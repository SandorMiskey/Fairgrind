package main

import "log/syslog"

const (

	// DB {{{

	// DB_CLEARING_BATCHES         string = "clearing_batches"
	// DB_CLEARING_BATCH_STATUS_ID string = "clearing_batch_status_id"
	// DB_CLEARING_BATCH_TYPE_ID   string = "clearing_batch_type_id"
	// DB_CLEARING_TASKS           string = "clearing_tasks"
	// DB_CLEARING_BATCH_ID     string = "clearing_batch_id"
	// DB_CLEARING_CLEARING_TASK_ID string = "clearing_clearing_task_id"
	DB_CLEARING_LEDGER_LABEL_TASK uint = 2

	/// }}}
	// ENV {{{

	ENV_PREFIX          string = "CLR_"
	ENV_CACHE_PASSWORD  string = ENV_PREFIX + "CACHE_PASSWORD"
	ENV_CACHE_HOST      string = ENV_PREFIX + "CACHE_HOST"
	ENV_DB_HOST         string = ENV_PREFIX + "DB_HOST"
	ENV_DB_PARAMS       string = ENV_PREFIX + "DB_PARAMS"
	ENV_DB_PASSWORD     string = ENV_PREFIX + "DB_PASSWORD"
	ENV_DB_PORT         string = ENV_PREFIX + "DB_PORT"
	ENV_DB_SCHEMA       string = ENV_PREFIX + "DB_SCHEMA"
	ENV_DB_USER         string = ENV_PREFIX + "DB_USER"
	ENV_LOG_SEVERITY    string = ENV_PREFIX + "LOG_SEVERITY"
	ENV_LOG_PREFIX      string = ENV_PREFIX + "LOG_PREFIX"
	ENV_MQ_EXCHANGE     string = ENV_PREFIX + "MQ_EXCHANGE"
	ENV_MQ_HOST         string = ENV_PREFIX + "MQ_HOST"
	ENV_MQ_PASSWORD     string = ENV_PREFIX + "MQ_PASSWORD"
	ENV_MQ_PORT         string = ENV_PREFIX + "MQ_PORT"
	ENV_MQ_ROUTING      string = ENV_PREFIX + "MQ_ROUTING"
	ENV_MQ_QUEUE        string = ENV_PREFIX + "MQ_QUEUE"
	ENV_MQ_USER         string = ENV_PREFIX + "MQ_USER"
	ENV_TICKER_INTERVAL string = ENV_PREFIX + "TICKER_INTERVAL"
	ENV_TICKER_MARKER   string = ENV_PREFIX + "TICKER_MARKER"

	/// }}}
	// ERR {{{

	ERR_DB_FAILED_TO_FETCH_BATCH_TYPE   string = "failed to fetch batch type"
	ERR_DB_FAILED_TO_FETCH_BATCH_STATUS string = "failed to fetch batch status"
	ERR_MQ_FAILED_TO_BIND               string = "failed to bind queue to exchange"
	ERR_MQ_FAILED_TO_CONNECT            string = "failed to connect to mq"
	ERR_MQ_FAILED_TO_CONSUME            string = "failed to register a consumer"
	ERR_MQ_FAILED_TO_DECLARE            string = "failed to declare a queue"
	ERR_MQ_FAILED_TO_PARSE              string = "failed to parse message"
	ERR_MQ_FAILED_TO_OPEN               string = "failed to open a channel"
	ERR_MQ_FAILED_TO_SET_QOS            string = "failed to set qoS"

	/// }}}
	// LOG {{{

	LOG_ERR    syslog.Priority = syslog.LOG_ERR
	LOG_NOTICE syslog.Priority = syslog.LOG_NOTICE
	LOG_INFO   syslog.Priority = syslog.LOG_INFO
	LOG_DEBUG  syslog.Priority = syslog.LOG_DEBUG
	LOG_EMERG  syslog.Priority = syslog.LOG_EMERG

	/// }}}
	// MSG {{{

	MSG_MSG_RECEIVED     string = "mq message received"
	MSG_MSG_PARSED       string = "message parsed successfully"
	MSG_ROUTING_DONE     string = "evaluation of the routing rules is completed"
	MSG_ROUTING_MATCH    string = "routing match"
	MSG_ROUTING_MISMATCH string = "no matching routing rule"
	MSG_SCHEMA_MISMATCH  string = "schema mismatch msg.Database != Env[ENV_DB_SCHEMA]"
	// MSG_SUBROUTING_DONE                 string = "sub-routing done"
	// MSG_SUBROUTING_BATCH_NONCREDITABLE  string = "sub-routing exception: batch is non-creditable"
	// MSG_SUBROUTING_MULTIPLIERS_MISMATCH string = "sub-routing exception: new multiplier is either less than or equal to the old one"
	// MSG_SUBROUTING_SKIPPING_FIELD string = "skipping field"
	MSG_TASK_BATCH_UPDATED string = "batch was updated later than the task"
	MSG_TASK_CLEARED       string = "task cleared"
	MSG_TASK_NONCREDITABLE string = "task is non-creditable"
	MSG_TASK_PROCESSING    string = "processing task"
	MSG_TASK_UNCLEARED     string = "uncleared task found"
	MSG_WAITING_FOREVER    string = "waiting forever"

	/// }}}

)
