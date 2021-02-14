package server

import (
	"encoding/json"
	"log"
	"os"
)

// Config ...
type Config struct {
	// Name - имя сервера
	Name string `json:name`
	// Limit - лимит сообщений
	Limit uint32 `json:limit`
	// Port - Порт сервера
	Port uint16 `json:port`
}

// ReadConfig ...
func ReadConfig() *[]Config {

	var config []Config

	file, err := os.Open("./server/config.json")

	if err != nil {
		log.Println("Не найден файл конфигурации")
		log.Fatal(err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Println("Не удалось прочитать файл конфигурации")
		log.Fatal(err)
	}

	return &config
}

// WriteConfig ...
func WriteConfig() {

	var config []Config

	config = append(config,
		Config{
			Name:  "Server1",
			Limit: 5,
		},
		Config{
			Name:  "Server2",
			Limit: 5,
		}, Config{
			Name:  "Server3",
			Limit: 5,
		})

	file, err := os.OpenFile("./server/config.json", os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	encoder := json.NewEncoder(file)
	encoder.Encode(&config)
}
