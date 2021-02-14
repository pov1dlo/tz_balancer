package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"guthub.com/povqdlo/tz_balancer/payload"
)

//CounterServer - Счетчик обработанных сообщений
var CounterServer uint64

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

			var payload []payload.Payload
			err := json.Unmarshal(data, &payload)
			if err != nil {
				log.Println("Ошибка чтения данных!", err)
			}

			CounterServer++

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

	var payload []payload.Payload
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
	defer s.mux.RUnlock()
	return s.blocked
}

// NewServer ...
func NewServer(c *Config) (s *Server) {

	s = &Server{}
	s.Name = c.Name
	s.MaxLimit = c.Limit
	s.Data = make(chan []byte, c.Limit)
	s.counter = 0
	s.blocked = false

	return
}
