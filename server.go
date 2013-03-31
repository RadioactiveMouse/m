package main

import (
	"bufio"
	"log"
	"os"
	"sync"
	"time"
)

// type to encapsulate the Server commands
type Server struct {
	name     string
	ch       chan Metric
	close    chan bool
	buf      []Metric
	interval int
	lk       sync.Mutex
}

// Make, initialise and run a metric server for the given name. Returns a type Server for the ability to defer closure
func NewMetricServer(name string, interval int) *Server {
	s := new(Server)
	s.name = name
	s.ch = make(chan Metric)
	s.close = make(chan bool)
	s.buf = make([]Metric, 0)
	if interval == 0 {
		s.interval = 60
	}
	go s.Run()
	return s
}

// run the server to listen for data across the channel
func (s *Server) Run() {
	defer s.Close()
	for {
		timer := time.Tick(60 * time.Second)
		select {
		case data := <-s.ch:
			ret := append(s.buf, data)
			if ret == nil {
				log.Printf("Metric %s was not written to the buffer\n", data.GetKey())
			}
		case <-timer:
			go s.Flush()
		case <-s.close:
			return
		}
	}
}

// flush the data in the buffer to the backend
func (s *Server) Flush() {
	var file *os.File
	var err error
	file, err = os.Open(s.name)
	if err != nil {
		// create the file as it's a new metric
		file, _ = os.Create(s.name)
	}
	defer file.Close()
	s.lk.Lock()
	defer s.lk.Unlock()
	w := bufio.NewWriter(file)
	// append to the file the contents of buf
	for _, metric := range s.buf {
		in, err := w.WriteString(metric.String())
		if err != nil || in == 0 {
			log.Println("Problem writing metric : ", metric)
		}
	}
}

// function to ensure correct and proper cleanup
func (s *Server) Close() {
	defer log.Printf("Metric server %s closed\n", s.name)
	close(s.ch)
	close(s.close)
}
