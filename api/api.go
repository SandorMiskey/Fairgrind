// region: packages

package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"log/syslog"
	"regexp"
	"strings"

	"utils"

	_ "api/swagger"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
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
		regexp.MustCompile("^/v1/swagger.*$"),
	}
)

const (
	LOG_ERR    syslog.Priority = syslog.LOG_ERR
	LOG_NOTICE syslog.Priority = syslog.LOG_NOTICE
	LOG_INFO   syslog.Priority = syslog.LOG_INFO
	LOG_DEBUG  syslog.Priority = syslog.LOG_DEBUG
	LOG_EMERG  syslog.Priority = syslog.LOG_EMERG

	ENV_PREFIX           string = "API_"
	ENV_APIKEY           string = ENV_PREFIX + "KEY"
	ENV_CACHE_PASSWORD   string = ENV_PREFIX + "CACHE_PASSWORD"
	ENV_CACHE_HOST       string = ENV_PREFIX + "CACHE_HOST"
	ENV_DB_HOST          string = ENV_PREFIX + "DB_HOST"
	ENV_DB_PARAMS        string = ENV_PREFIX + "DB_PARAMS"
	ENV_DB_PASSWORD      string = ENV_PREFIX + "DB_PASSWORD"
	ENV_DB_PORT          string = ENV_PREFIX + "DB_PORT"
	ENV_DB_SCHEMA        string = ENV_PREFIX + "DB_SCHEMA"
	ENV_DB_USER          string = ENV_PREFIX + "DB_USER"
	ENV_FIBER_ADDRESS    string = ENV_PREFIX + "FIBER_ADDRESS"
	ENV_FIBER_LOGFORMAT  string = ENV_PREFIX + "FIBER_LOGFORMAT"
	ENV_FIBER_TIMEFORMAT string = ENV_PREFIX + "FIBER_TIMEFORMAT"
	ENV_FIBER_METRICS    string = ENV_PREFIX + "FIBER_METRICS"
	ENV_LOG_SEVERITY     string = ENV_PREFIX + "LOG_SEVERITY"
	ENV_LOG_PREFIX       string = ENV_PREFIX + "LOG_PREFIX"
	ENV_SWAGGER_BASEPATH string = ENV_PREFIX + "SWAGGER_BASEPATH"
	ENV_SWAGGER_FILEPATH string = ENV_PREFIX + "SWAGGER_FILEPATH"
	ENV_SWAGGER_PATH     string = ENV_PREFIX + "SWAGGER_PATH"
	ENV_SWAGGER_TITLE    string = ENV_PREFIX + "SWAGGER_TITLE"
)

// endregion

func main() {

	// region: read .env

	// Env, err = godotenv.Read()
	// utils.Panic(err)
	Env = utils.GetEnv()

	// endregion: read .env
	// region: logger

	utils.GetLogger(Env[ENV_LOG_SEVERITY], Env[ENV_LOG_PREFIX])
	defer utils.Logger.Close()
	Logger = utils.Logger.Out

	// endregion
	// region: db

	// TODO: use "database/sql" for connection pool and fine tunning

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", Env[ENV_DB_USER], Env[ENV_DB_PASSWORD], Env[ENV_DB_HOST], Env[ENV_DB_PORT], Env[ENV_DB_SCHEMA], Env[ENV_DB_PARAMS])
	Logger(LOG_DEBUG, dsn)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	utils.Panic(err)

	// endregion
	// region: redis

	Cache = redis.NewClient(&redis.Options{
		Addr:     Env[ENV_CACHE_HOST],
		Password: Env[ENV_CACHE_PASSWORD],
		DB:       0,
		Protocol: 3,
	})
	Logger(LOG_DEBUG, Cache)

	// endregion: redis
	// region: fiber

	// region: new app w/ logger

	api := fiber.New(fiber.Config{
		Prefork: false,
	})

	api.Use(logger.New(logger.Config{
		Format:     Env[ENV_FIBER_LOGFORMAT] + "\n",
		TimeFormat: Env[ENV_FIBER_TIMEFORMAT],
	}))

	// endregion
	// region: authentication

	api.Use(keyauth.New(keyauth.Config{
		Next: func(c *fiber.Ctx) bool {
			originalURL := strings.ToLower(c.OriginalURL())

			for _, pattern := range unprotectedURLs {
				if pattern.MatchString(originalURL) {
					Logger(LOG_NOTICE, originalURL, "unauthenticated")
					return true
				}
			}
			return false
		},
		Validator: func(c *fiber.Ctx, key string) (bool, error) {
			// Logger(LOG_DEBUG, "request headers", c.GetReqHeaders())
			hashedAPIKey := sha256.Sum256([]byte(Env[ENV_APIKEY]))
			hashedKey := sha256.Sum256([]byte(key))
			if subtle.ConstantTimeCompare(hashedAPIKey[:], hashedKey[:]) == 1 {
				return true, nil
			}
			Logger(LOG_DEBUG, c.OriginalURL(), "authentication failure")
			return false, keyauth.ErrMissingOrMalformedAPIKey
		},
	}))

	// endregion
	// region: other middleware

	api.Use(cors.New())
	api.Use(recover.New())
	api.Use(requestid.New())
	api.Use(healthcheck.New())
	api.Get("Env[ENV_FIBER_METRICS]", monitor.New())

	// endregion
	// region: swagger

	//	@securityDefinitions.apikey	ApiKeyAuth
	//	@in							header
	//	@name						Authorization

	//	@BasePath	/v1

	api.Use(swagger.New(swagger.Config{
		BasePath: Env[ENV_SWAGGER_BASEPATH],
		FilePath: Env[ENV_SWAGGER_FILEPATH],
		Path:     Env[ENV_SWAGGER_PATH],
		Title:    Env[ENV_SWAGGER_TITLE],
	}))

	// endregion
	// region: routes

	v1 := api.Group(Env[ENV_SWAGGER_BASEPATH])

	// region: batches

	v1_batches := v1.Group("/batches")
	v1_batches.Get("/types", v1_batches_types_get)
	v1_batches.Get("/statuses", v1_batches_statuses_get)

	// endregion
	// region: ledger

	v1_ledger := v1.Group("/ledger")
	v1_ledger.Get("/labels", v1_ledger_labels_get)
	v1_ledger.Get("/statuses", v1_ledger_statuses_get)

	v1_ledger.Get("", v1_ledger_get)
	v1_ledger.Delete("", v1_ledger_delete)
	// TODO
	// * credit (post)
	// * swap (patch)
	// * withdraw (delete)
	// * transfer (put)

	// endregion
	// region: tasks

	v1_tasks := v1.Group("/tasks")
	v1_tasks.Get("", v1_tasks_get)
	v1_tasks.Post("", v1_tasks_post)
	v1_tasks.Delete("/fees", v1_tasks_fees_delete)
	v1_tasks.Get("/fees", v1_tasks_fees_get)
	v1_tasks.Post("/fees", v1_tasks_fees_post)
	v1_tasks.Get("/types", v1_tasks_types_get)
	v1_tasks.Get("/statuses", v1_tasks_statuses_get)

	// endregion
	// region: tokens

	v1_tokens := v1.Group("/tokens")
	v1_tokens.Get("", v1_tokens_get)
	v1_tokens.Get("/types", v1_tokens_types_get)

	// endregion
	// region: wallets

	v1_wallets := v1.Group("/wallets")
	v1_wallets.Get("/summed", v1_wallets_summed_get)
	v1_wallets.Get("/detailed", v1_wallets_detailed_get)

	// endregion

	api.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	// endregion
	// region: listen

	err := api.Listen(Env[ENV_FIBER_ADDRESS])
	utils.Panic(err)

	// endregion

	// endregion

}
