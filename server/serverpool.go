package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
)

// ServerPool ...
type ServerPool []*Server

// LoadBalance ...
func (s *ServerPool) LoadBalance(pool *ServerPool, w http.ResponseWriter, r *http.Request) {

	var data []byte
	data, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Ошибка, %s", err)
	}
	defer r.Body.Close()

	sort.Sort(sort.Reverse(pool))

	for _, server := range *pool {
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
