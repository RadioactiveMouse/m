package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/RadioactiveMouse/rgo"
)

// type to encapsulate the Server commands
type Server struct {
	name     string
	ch       chan Metric
	close    chan bool
	lk       sync.Mutex
	buf      []Metric
	interval int
	bucket	rgo.Bucket
	errors []error
}

func (s *Server) log(content interface{}) {
	log.Println(fmt.Sprintf("[%s] : %v", s.name, content))
}

// Make, initialise and run a metric server for the given name. Returns a type Server for the ability to defer closure
func NewMetricServer(name string, interval int) *Server {
	s := new(Server)
	s.name = name
	s.ch = make(chan Metric)
	s.close = make(chan bool)
	s.buf = make([]Metric, 0)
	s.errors = make([]error, 0)
	s.bucket = rgo.HTTPClient().Bucket(name)
	if interval == 0 {
		s.interval = 60
	}
	// check to make sure the server connection exists
	err := rgo.Ping("192.168.1.10")
	if err != nil {
		log.Fatal("Could not contact the riak cluster.")
	}
	go s.Run()
	return s
}

// run the server to listen for data across the channel
func (s *Server) Run() {
	defer s.Close()
	timer := time.Tick(s.interval * time.Second)
	for {
		select {
		case data := <-s.ch:
			s.buf = append(s.buf, data)
			s.log("Metric written to buffer")
		case <-timer:
			go s.Flush()
		case <-s.close:
			return
		}
	}
}

// flush the data in the buffer to the backend
func (s *Server) Flush() {
	s.lk.Lock()
	if s.buf == nil {
		s.lk.Unlock()
		return
	} else {
		// copy the buffer and reset the server buffer
		buffer := make([]Metric,len(s.buf))
		buffer = s.buf
		s.buf = nil
		s.lk.Unlock()
	}
	for _, metric := range buffer {
		object := s.bucket.Data()
		object.Key := time.Now().String()
		object.Value := []byte("Test") //metric.String()
		_, err := object.Store()
		if err != nil {
			s.errors = append(errors, err)
		}
	}
}

// function to ensure correct and proper cleanup
func (s *Server) Close() {
	defer log.Printf("Metric server [%s] closed\n", s.name)
	close(s.ch)
	close(s.close)
}
