package config

import "github.com/joho/godotenv"

func MustInit(path string) {
	if err := godotenv.Load(path); err != nil {
		panic(err)
	}
}
