package files

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/markbates/pkger"

	log "github.com/sirupsen/logrus"
)

func init() {
	env, err := pkger.Open("/.env")
	if err != nil {
		log.Fatal(err)
	}

	envMap, err := godotenv.Parse(env)
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range envMap {
		os.Setenv(key, value)
	}
}
