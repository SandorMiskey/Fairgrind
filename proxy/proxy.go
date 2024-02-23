// region: packages

package main

import (
	"log/syslog"
	"regexp"
	"utils"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// endregion: packages
// region: globals

var (
	err error

	DB     *gorm.DB
	Env    map[string]string
	Cache  *redis.Client
	Logger func(s ...interface{}) *[]error

	unprotectedURLs = []*regexp.Regexp{
		regexp.MustCompile("^/api/v1/swagger$"),
	}
)

const (
	LOG_ERR    syslog.Priority = syslog.LOG_ERR
	LOG_NOTICE syslog.Priority = syslog.LOG_NOTICE
	LOG_INFO   syslog.Priority = syslog.LOG_INFO
	LOG_DEBUG  syslog.Priority = syslog.LOG_DEBUG
	LOG_EMERG  syslog.Priority = syslog.LOG_EMERG

	ENV_PREFIX       string = "PROXY_"
	ENV_APIKEY       string = ENV_PREFIX + "APIKEY"
	ENV_LOG_SEVERITY string = ENV_PREFIX + "LOG_SEVERITY"
	ENV_LOG_PREFIX   string = ENV_PREFIX + "LOG_PREFIX"
	// ENV_CACHE_PASSWORD   string = ENV_PREFIX + "CACHE_PASSWORD"
	// ENV_CACHE_HOST       string = ENV_PREFIX + "CACHE_HOST"
	// ENV_DB_HOST          string = ENV_PREFIX + "DB_HOST"
	// ENV_DB_PARAMS        string = ENV_PREFIX + "DB_PARAMS"
	// ENV_DB_PASSWORD      string = ENV_PREFIX + "DB_PASSWORD"
	// ENV_DB_PORT          string = ENV_PREFIX + "DB_PORT"
	// ENV_DB_SCHEMA        string = ENV_PREFIX + "DB_SCHEMA"
	// ENV_DB_USER          string = ENV_PREFIX + "DB_USER"
	// ENV_FIBER_ADDRESS    string = ENV_PREFIX + "FIBER_ADDRESS"
	// ENV_FIBER_LOGFORMAT  string = ENV_PREFIX + "FIBER_LOGFORMAT"
	// ENV_FIBER_TIMEFORMAT string = ENV_PREFIX + "FIBER_TIMEFORMAT"
	// ENV_SWAGGER_BASEPATH string = ENV_PREFIX + "SWAGGER_BASEPATH"
	// ENV_SWAGGER_FILEPATH string = ENV_PREFIX + "SWAGGER_FILEPATH"
	// ENV_SWAGGER_PATH     string = ENV_PREFIX + "SWAGGER_PATH"
	// ENV_SWAGGER_TITLE    string = ENV_PREFIX + "SWAGGER_TITLE"
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
	Logger = utils.Logger.Out

	// endregion

}
