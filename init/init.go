// region: packages

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"models"
	"utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// endregion: packages
// region: const

var (
	Action string = "templates"
)

const (
	DB_HOST     string = "DB_HOST"
	DB_PARAMS   string = "DB_PARAMS"
	DB_PASSWORD string = "DB_PASSWORD"
	DB_PORT     string = "DB_PORT"
	DB_SCHEMA   string = "DB_SCHEMA"
	DB_USER     string = "DB_USER"

	PATH_WORKBENCH string = "PATH_WORKBENCH"
)

// endregion

func main() {

	// region: read .env

	// env, err := godotenv.Read()
	// if err != nil {
	// 	utils.Panic(err)
	// }
	env := utils.GetEnv()

	// endregion: read .env
	// region: default action

	if len(os.Args) > 1 {
		Action = os.Args[1]
	}

	// endregion: default action
	// region: templates

	if Action == "templates" {
		err := filepath.Walk(env["PATH_TEMPLATES"], func(path string, info os.FileInfo, err error) error {

			// region: can I walk here

			if err != nil {
				fmt.Println(err) // can't walk here,
				return nil       // but continue walking elsewhere
			}

			// endregion
			// region: target path

			var builder strings.Builder
			var target string

			builder.WriteString(env["PATH_WORKBENCH"])
			builder.WriteString(path[len(env["PATH_TEMPLATES"]):])

			target = builder.String()

			// endregion
			// region: apply template

			if info.IsDir() {
				_, err := os.Stat(target)
				if os.IsNotExist(err) {
					err = os.MkdirAll(target, 0755)
					if err != nil {
						utils.Panic(err)
					}
					fmt.Println("directory created:", path, " -> ", target)
					return nil
				}
			} else {

				// region: read file

				content, err := os.ReadFile(path)
				if err != nil {
					utils.Panic(err)
				}
				contentType := http.DetectContentType(content)

				// endregion
				// region: process

				if strings.HasPrefix(contentType, "text/plain") {
					temp, err := template.New("template").Parse(string(content))
					if err != nil {
						utils.Panic(err)
					}

					outputFile, err := os.Create(target)
					if err != nil {
						return err
					}
					defer outputFile.Close()

					err = temp.Execute(outputFile, env)
					if err != nil {
						utils.Panic(err)
					}

					fmt.Printf("template applied: %s -> %s (%s)\n", path, target, contentType)

					return nil
				}

				// endregion
				// region: copy

				sourceFile, err := os.Open(path)
				if err != nil {
					return err
				}
				defer sourceFile.Close()

				destinationFile, err := os.Create(target)
				if err != nil {
					return err
				}
				defer destinationFile.Close()

				_, err = io.Copy(destinationFile, sourceFile)
				if err != nil {
					return err
				}

				fmt.Printf("file copied: %s -> %s (%s)\n", path, target, contentType)

				// endregion

			}

			return nil

			// endregion

		})
		if err != nil {
			fmt.Printf("error walking the path %v: %v\n", env["PATH_TEMPLATES"], err)
		}
	}

	// endregion: templates
	// region: gorm init

	if Action == "gorm" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", env[DB_USER], env[DB_PASSWORD], env[DB_HOST], env[DB_PORT], env[DB_SCHEMA], env[DB_PARAMS])
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			utils.Panic(err)
		}

		db.AutoMigrate(&models.ClearingBatchType{})
		db.AutoMigrate(&models.ClearingBatchStatus{})
		db.AutoMigrate(&models.ClearingBatch{})

		db.AutoMigrate(&models.ClearingLedgerStatus{})
		db.AutoMigrate(&models.ClearingLedgerLabel{})
		db.AutoMigrate(&models.ClearingLedger{})

		db.AutoMigrate(&models.ClearingTaskStatus{})
		db.AutoMigrate(&models.ClearingTaskType{})
		db.AutoMigrate(&models.ClearingTaskFee{})
		db.AutoMigrate(&models.ClearingTask{})

		db.AutoMigrate(&models.ClearingTokenType{})
		db.AutoMigrate(&models.ClearingToken{})
	}

	// endregion

}
