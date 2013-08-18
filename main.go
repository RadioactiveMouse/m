package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	port        = flag.Int("port", 2013, "port to host the m server on")
	buf         = make(map[string]chan Metric)
	parseErrors = 0
)

const interval = 60

func main() {
	flag.Parse()
	// spin up a connection to accept metrics being sent to the server
	c, err := net.Listen("tcp", ":2013")
	if err != nil {
		log.Fatal("Couldn't open a connection to accept Metrics")
	}
	for {
		conn, er := c.Accept()
		if er != nil {
			log.Println("Couldn't accept a connection")
		}
		// read from the connection and pipe it to the correct channel
		go connection(conn)
	}
}

func connection(conn net.Conn) {
	buffer := make([]byte, 100) // maybe change this value later on in case metrics get large?
	for {
		n, readError := conn.Read(buffer)
		if readError != nil {
			if readError == io.EOF {
				return
			}
			log.Println(readError)
			return
		}
		defer conn.Close()
		met, parseError := parseincoming(buffer[:n])
		if parseError != nil {
			log.Println(parseError)
			parseErrors = parseErrors + 1
		} else {
			key := met.GetKey()
			// check if met.key is in the map of current servers
			if _, ok := buf[key]; ok {
				// key exists so send value to waiting goroutine
				buf[key] <- met
			} else {
				// make new metrics server
				ms := NewMetricServer(key, interval)
				buf[key] = ms.ch
				buf[key] <- met
			}
		}
	}
}

// parse the incoming content into a metric type
func parseincoming(content []byte) (Metric, error) {
	// try to unmarshal
	var m Metric
	err := json.Unmarshal(content, &m)
	return m, err
}
