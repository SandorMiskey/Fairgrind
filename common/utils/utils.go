// region: packages

package utils

import (
	"fmt"
	"log/syslog"
	"models"
	"os"
	"strconv"
	"strings"

	"github.com/SandorMiskey/TEx-kit/log"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// endregion
// region: globals

var (
	Logger      *log.Logger
	LogPrefix   string          = "-> "
	LogSeverity syslog.Priority = log.LOG_INFO
)

// endregion: globals

func Panic(err error, s ...string) {
	if err != nil {
		s = append([]string{err.Error()}, s...)
		msg := strings.Join(s, " -> ")
		if Logger != nil {
			Logger.Out(syslog.LOG_EMERG, msg)
		}
		fmt.Fprintln(os.Stderr, msg)
		os.Exit(1)
	}
}

func GetEnv(dotenv ...string) map[string]string {
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
}

func GetLogger(p ...string) *log.Logger {

	// region: settings

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

	// endregion
	// region: mute

	if severity < syslog.LOG_INFO {
		log.ChDefaults.Welcome = nil
		log.ChDefaults.Bye = nil
	} else {
		welcome := fmt.Sprintf("%s (level: %v)", *log.ChDefaults.Welcome, severity)
		log.ChDefaults.Welcome = &welcome
	}

	// endregion
	Logger = log.NewLogger()
	_, _ = Logger.NewCh(log.ChConfig{Severity: &severity, Prefix: &prefix})
	return Logger
}

func GetResponse(c *fiber.Ctx) models.ApiResponse {
	meta := make(models.ApiResponseMeta)
	meta["queries"] = c.Queries()
	meta["body_raw"] = c.BodyRaw()
	// meta["req_headers"] = c.GetReqHeaders()

	return models.ApiResponse{
		Meta:    meta,
		Success: false,
	}
}
