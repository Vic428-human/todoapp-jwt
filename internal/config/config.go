package config

import (
	"log"
	"os" //retrieve environment variables

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Port        string
}

/*
make api call to database
*/
// pointer : point to the address or instance in memory
func Load() (*Config, error) {
	// loads .env from the current directory
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	/* *pointer (point to a struct) = &reference{} (reference to a struct as value) => https://stackoverflow.com/questions/47296325/passing-by-reference-and-value-in-go-to-functions
		func someFunc(x *int) {
	    *x = 2 // Whatever variable caller passed in will now be 2
	    y := 7
	    x = &y // has no impact on the caller because we overwrote the pointer value!
	}
	*/
	var config *Config = &Config{
		DatabaseURL: os.Getenv("POSTGRES_URL"),
		Port:        os.Getenv("POSTGRES_PORT"),
	}
	return config, nil
}
