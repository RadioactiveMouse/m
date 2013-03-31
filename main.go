package main

import (
	"flag"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
	"errors"
)

var (
	port   = flag.Int("port", 2013, "port to host the m server on")
	config = flag.Bool("config", false, "switch on/off debugging information")
	buf    = make(map[string]chan Metric)
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
	spinUp()
	conn, er := c.Accept()
	if er != nil {
		log.Fatal("Couldn't accept a connection")
	}
	// read from the connection and pipe it to the correct channel
	buffer := make([]byte, 100)
	for {
		_, readError := conn.Read(buffer)
		if readError != nil {
			log.Fatal(readError)
		}
		defer c.Close()
		met, parseError := parse(string(buffer))
		if parseError != nil {
			log.Println("Error parsing a metric")
			parseErrors = parseErrors + 1
		} else {
			key := met.GetKey()
			// check if met.key is in the map of current servers
			if _, ok := buf[key]; ok {
				// key exists
				buf[key] <- met
			} else {
				// make new metrics server
				ms := NewMetricServer(key, interval)
				ms.ch <- met
				buf[key] = ms.ch
			}
		}
	}
}

// look in the directory and start any old metric servers
func spinUp() {
	return
}

// parse the data in form key|value
func parse(s string) (Metric, error) {
	data := strings.Split(s, "|")
	in, er := strconv.ParseFloat(data[1], 32)
	if er == nil {
		c := new(Counter)
		c.SetKey(data[0])
		c.SetValue(in)
		return c, nil
	}
	ti, err := time.Parse(time.RFC1123Z, data[1])
	if err == nil {
		t := new(TimeSeries)
		t.SetKey(data[0])
		t.SetValue(ti)
		return t, nil
	}
	return nil, errors.New("Non recognised type wasn't parsed")
}
