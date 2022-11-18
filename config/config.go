package config

import "os"

var (
	DB_FILEPATH      string
	AES_IV           string // must be of length 32
	AES_SPLICE_INDEX string // must be a string parsable to >=0 and <= 31
)

/* Sets environment variables ued by program. Will be different for integration
tests than in production */
func SetupEnv(isProduction bool) {
	if isProduction {
		DB_FILEPATH = "./data/blog.db"
		AES_IV = "snooping6is9bad0"
		AES_SPLICE_INDEX = "28"
	} else {
		DB_FILEPATH = "./data/blog_test.db"
		AES_IV = "snooping6is9bad0"
		AES_SPLICE_INDEX = "28"
	}
	os.Setenv("DB_FILEPATH", DB_FILEPATH)
	os.Setenv("AES_IV", AES_IV)
	os.Setenv("AES_SPLICE_INDEX", AES_SPLICE_INDEX)
}
