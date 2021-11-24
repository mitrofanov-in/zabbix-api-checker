package getenv

import (
	"os"
	// Import godotenv
	"github.com/joho/godotenv"
)



// use godot package to load/read the .env file and
// return the value of the key
func GoDotEnvVariable(key string) string {

	// load .env file
godotenv.Load("C:/Users/user/Documents/.env")
	return os.Getenv(key)
}
