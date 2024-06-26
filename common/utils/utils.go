// region: packages {{{

package utils

import (
	"encoding/json"
	"fmt"
	"log/syslog"
	"models"
	"os"
	"strconv"
	"strings"

	"github.com/SandorMiskey/TEx-kit/log"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// endregion }}}
// region: globals {{{

var (
	Logger      *log.Logger
	LogPrefix   string          = "-> "
	LogSeverity syslog.Priority = log.LOG_INFO
)

// endregion: globals }}}

func Contains[T comparable](slice []T, val T) bool { // {{{
	// Contains checks if a slice contains a given value.
	// The function uses generics to work with any comparable type.
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
} // }}}

func GetEnv(dotenv ...string) map[string]string { // {{{
	err := godotenv.Load(dotenv...)
	if err != nil {
		if Logger != nil {
			Logger.Out(syslog.LOG_EMERG, err)
		}
		fmt.Fprintln(os.Stderr, err.Error())
	}

	envVars := make(map[string]string)
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			envVars[pair[0]] = pair[1]
		}
	}

	return envVars
} // }}}

func GetLogger(p ...string) *log.Logger { // {{{

	// region: settings {{{

	severity := syslog.Priority(LogSeverity)
	prefix := LogPrefix
	if len(p) > 0 {
		s, err := strconv.ParseInt(p[0], 10, 32)
		Panic(err)
		severity = syslog.Priority(s)
	}
	if len(p) > 1 {
		prefix = p[1]
	}

	// endregion }}}
	// region: mute {{{

	if severity < syslog.LOG_INFO {
		log.ChDefaults.Welcome = nil
		log.ChDefaults.Bye = nil
	} else {
		welcome := fmt.Sprintf("%s (level: %v)", *log.ChDefaults.Welcome, severity)
		log.ChDefaults.Welcome = &welcome
	}

	// endregion }}}

	Logger = log.NewLogger()
	_, _ = Logger.NewCh(log.ChConfig{Severity: &severity, Prefix: &prefix})
	return Logger
} // }}}

func GetResponse(c *fiber.Ctx) models.ApiResponse { // {{{
	return models.ApiResponse{
		Request: models.ApiRequest{
			BodyRaw: string(c.BodyRaw()),
			Queries: c.Queries(),
			// ReqHeaders: c.GetReqHeaders(),
		},
		Success: false,
	}
} // }}}

func Panic(err error, s ...string) { // {{{
	if err != nil {
		s = append([]string{err.Error()}, s...)
		msg := strings.Join(s, " -> ")
		if Logger != nil {
			Logger.Out(syslog.LOG_EMERG, msg)
		}
		fmt.Fprintln(os.Stderr, msg)
		os.Exit(1)
	}
} // }}}

func Paginate(page, size int) func(db *gorm.DB) *gorm.DB { // {{{
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case size > 100:
			size = 100
		case size <= 0:
			size = 10
		}

		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
} // }}}

func JsonPP(msg interface{}) string { // {{{
	var bytes []byte
	var err error

	bytes, err = json.MarshalIndent(msg, "", "    ")
	if err != nil {
		bytes = []byte(fmt.Sprintf("%v (%s)", msg, err))
	}

	return string(bytes)
} // }}}
