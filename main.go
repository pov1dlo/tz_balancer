package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
)

//counterIn - Счетчик входящих сообщений
var counterIn uint64

//counterOut - Счетчик обработанных сообщений
var counterOut uint64

// Payload ...
type Payload struct {
	Price    int `json:price`
	Quantity int `json:quantity`
	Amount   int `json:amount`
	Object   int `json:object`
	Method   int `json:method`
}

// Config ...
type Config struct {
	//Name - имя сервера
	Name string `json:name`
	//Limit - лимит сообщений
	Limit uint32 `json:limit`
	//Port - Порт сервера
	Port uint16 `json:port`
}

// ServerPool ...
type ServerPool []*Server

// LoadBalance ...
func (s *ServerPool) LoadBalance(w http.ResponseWriter, r *http.Request) {

	var data []byte
	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Ошибка, %s", err)
	}
	defer r.Body.Close()

	sort.Sort(sort.Reverse(Pool))

	for _, server := range Pool {
		if !server.IsBlocked() {
			//go server.ProcessedData(w, data)
			server.Data <- data
			fmt.Fprint(w, fmt.Sprintf("Запрос передан на %s", server.Name))
			break
		}
	}
}

func (s ServerPool) Len() int {
	return len(s)
}

func (s ServerPool) Less(i, j int) bool {
	return s[i].counter > s[j].counter
}

func (s ServerPool) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

var config []Config

// Pool ...
var Pool ServerPool

// Limitter ...
type Limitter struct {
	MaxLimit uint32
	counter  uint32
	mux      sync.RWMutex
}

// Server ...
type Server struct {
	Name    string
	Data    chan []byte
	blocked bool
	Limitter
}

//Run ...
func Run(s *Server) {

	go func() {
		for data := range s.Data {

			log.Printf("Обработка пакета данных на %s\n", s.Name)

			s.mux.Lock()
			s.counter++
			if s.MaxLimit == s.counter {
				s.blocked = true
			}
			s.mux.Unlock()

			var payload []Payload
			err := json.Unmarshal(data, &payload)
			if err != nil {
				log.Println("Ошибка чтения данных!", err)
			}

			counterOut++

			s.mux.Lock()
			s.counter--
			if s.counter == 0 {
				s.blocked = false
			}
			s.mux.Unlock()
		}
	}()
}

//ProcessedData ... рудимент
func (s *Server) ProcessedData(w http.ResponseWriter, data []byte) {

	log.Printf("Обработка пакета данных на %s", s.Name)
	s.mux.Lock()
	s.counter++
	if s.MaxLimit == s.counter {
		s.blocked = true
	}
	s.mux.Unlock()

	var payload []Payload
	err := json.Unmarshal(data, &payload)
	if err != nil {
		log.Println("Ошибка чтения данных!", err)
	}

	fmt.Fprint(w, fmt.Sprintf("%s Успех!", s.Name))
	s.mux.Lock()
	s.counter--
	if s.MaxLimit != s.counter {
		s.blocked = false
	}
	s.mux.Unlock()

}

// IsBlocked ...
func (s *Server) IsBlocked() bool {
	s.mux.RLock()
	s.mux.RUnlock()
	return s.blocked
}

func readConfig() []Config {

	file, err := os.Open("config.json")

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

	return config
}

func writeConfig() {

	config = append(config,
		Config{
			Name:  "Server1",
			Limit: 1,
		},
		Config{
			Name:  "Server2",
			Limit: 1,
		}, Config{
			Name:  "Server3",
			Limit: 1,
		})

	file, err := os.OpenFile("config.json", os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	encoder := json.NewEncoder(file)
	encoder.Encode(&config)
}

// NewServer ...
func NewServer(c Config) (s *Server) {

	s = &Server{}
	s.Name = c.Name
	s.MaxLimit = c.Limit
	s.Data = make(chan []byte, c.Limit)
	s.counter = 0
	s.blocked = false

	return
}

func init() {

	for _, c := range readConfig() {

		newServer := NewServer(c)
		Run(newServer)
		Pool = append(Pool, newServer)

	}
}

func main() {

	fs := http.FileServer(http.Dir("./client"))
	http.Handle("/", fs)

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			counterIn++
			Pool.LoadBalance(w, r)
		}
	})

	lbPort := "8080"
	go http.ListenAndServe(fmt.Sprintf(":%s", lbPort), nil)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	fmt.Printf("Количество входящих: %d\tКоличество обработанных: %d", counterIn, counterOut)

}
