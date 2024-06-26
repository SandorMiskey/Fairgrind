module proxy

go 1.21.9

replace models => ../common/models

replace utils => ../common/utils

replace mq => ../common/mq

require (
	github.com/joho/godotenv v1.5.1
	github.com/redis/go-redis/v9 v9.5.1
	gorm.io/gorm v1.25.9
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
)
