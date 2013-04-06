package main

import (
	"bufio"
	"log"
	"os"
	"sync"
	"time"
	"fmt"
)

// type to encapsulate the Server commands
type Server struct {
	name     string
	ch       chan Metric
	close    chan bool
	lk       sync.Mutex
	buf      []Metric
	interval int
}

func (s *Server) log(content interface{}) {
	log.Println(fmt.Sprintf("[%s] : %v",s.name,content))
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
	// check to make sure the data directory exists and if not create it
	dirErr := os.Mkdir("data",0700)
	if dirErr != nil && !os.IsExist(dirErr) {
		log.Fatal("Error making the directory to store the metrics")
	}
	go s.Run()
	return s
}

// run the server to listen for data across the channel
func (s *Server) Run() {
	defer s.Close()
	timer := time.Tick(30 * time.Second)
	for {
		select {
		case data := <-s.ch:
			ret := append(s.buf, data)
			s.buf = ret
			if ret == nil {
				s.log("Metric was not written to the buffer")
			} else {
				s.log("Metric written to buffer")
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
	s.lk.Lock()
	defer s.lk.Unlock()
	if s.buf == nil {
		return
	}
	file, err := os.OpenFile(fmt.Sprintf("data/%s",s.name),os.O_APPEND|os.O_CREATE|os.O_RDWR,0666)
	if err != nil {
		s.log(err)
		return
	}
	defer file.Close()
	w := bufio.NewWriter(file)
	// append to the file the contents of buf
	for _, metric := range s.buf {
		in, err := w.WriteString(metric.String()+"\n")
		if err != nil || in == 0 {
			s.log(fmt.Sprintf("Problem writing metric : %s", metric))
		}
	}
	flushErr := w.Flush()
	if flushErr != nil {
		s.log(flushErr)
	} else {
		s.log("Metrics written to file")
		// reset the buffer
		s.buf = nil
	}
}

// function to ensure correct and proper cleanup
func (s *Server) Close() {
	defer log.Printf("Metric server %s closed\n", s.name)
	close(s.ch)
	close(s.close)
}
