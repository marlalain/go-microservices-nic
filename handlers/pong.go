package handlers

import (
	"log"
	"net/http"
)

type Pong struct {
	l *log.Logger
}

func NewPong(l *log.Logger) *Pong {
	return &Pong{l}
}

func (p *Pong) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	p.l.Println("server: Ping? Pong!")
	_, err := w.Write([]byte("Pong"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
