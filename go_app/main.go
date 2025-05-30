package main

import (
	"log"
	"os"
)

func main() {
	env_validate()
	Log_init()
	DB_init()
	data_mutations()

	Put("cat", "meow")
	Put("dog", "woof")
	Put("fish", "blub")

	val := Get("fish")
	log.Printf("retrieved_val: %v", val)

	Dump()

	select {}
}

func env_validate() {
	requiredVars := []string{
		"DB_HOST",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"APP_LOG_DIR",
	}

	missing := false
	for _, v := range requiredVars {
		if val := os.Getenv(v); val == "" {
			log.Printf("ERROR: Environment variable %s is not set or empty", v)
			missing = true
		}
	}
	if missing {
		log.Fatal("One or more required env vars are missing, quit.")
	}
}
