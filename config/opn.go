package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/omise/omise-go"
)

func InitOmise() *omise.Client {
	err := godotenv.Load()
	if err!= nil {
        log.Fatal("Error loading.env file")
    }

	var(
		OmisePublicKey = os.Getenv("OmisePublicKey")
		OmiseSecretKey = os.Getenv("OmiseSecretKey")
	)

	client, e := omise.NewClient(OmisePublicKey, OmiseSecretKey)
	if e != nil {
		log.Fatal(e);
	}

	return client;
}