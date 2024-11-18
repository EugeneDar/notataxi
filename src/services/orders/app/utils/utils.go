package utils

import (
	"log"
	"os"

	"github.com/google/uuid"
)

func GetenvSafe(name string) string {
	env, doesItExist := os.LookupEnv(name)
	if !doesItExist {
		log.Fatalf("There is no %s in the environment variables", name)
	}
	return env
}

func GenerateUUID() string {
	return uuid.New().String()
}
