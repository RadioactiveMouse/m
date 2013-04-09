package main

import (
	"flag"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
	"errors"
	"io"
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
	for {
		conn, er := c.Accept()
		if er != nil {
			log.Fatal("Couldn't accept a connection")
		}
		// read from the connection and pipe it to the correct channel
		go connection(conn)
	}
}

func connection(conn net.Conn) {
	buffer := make([]byte, 100)
	for {
		n, readError := conn.Read(buffer)
		if readError != nil {
			if readError == io.EOF {
				return
			}
			log.Fatal(readError)
		}
		defer conn.Close()
		met, parseError := parse(string(buffer[:n]))
		if parseError != nil {
			log.Println(parseError)
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
	// trim the newline off the end of the value
	data[1] = strings.TrimSpace(data[1])
	fl, er := strconv.ParseFloat(data[1], 32)
	if er == nil {
		c := new(Counter)
		c.SetKey(data[0])
		c.SetValue(fl)
		return c, nil
	}
	/* in, intError := strconv.ParseInt(data[1],10,32)
	if intError == nil {
		c := new(Counter)
		c.SetKey(data[0])
		c.SetValue(float64(in))
		log.Println(in)
		return c, nil
	} */
	ti, err := time.Parse(time.RFC1123Z, data[1])
	if err == nil {
		t := new(TimeSeries)
		t.SetKey(data[0])
		t.SetValue(ti)
		return t, nil
	}
	return nil, errors.New("Non recognised type wasn't parsed")
}
