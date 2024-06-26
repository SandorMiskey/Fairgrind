module clearing

go 1.21.9

replace models => ../common/models

replace utils => ../common/utils

require (
	github.com/rabbitmq/amqp091-go v1.9.0
	gorm.io/driver/mysql v1.5.6
	gorm.io/gorm v1.25.9
	models v0.0.0-00010101000000-000000000000
	utils v0.0.0-00010101000000-000000000000
)

require (
	github.com/SandorMiskey/TEx-kit v0.0.1 // indirect
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/gofiber/fiber/v2 v2.52.4 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.17.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.51.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
)
