package main

import (
	"time"
	"fmt"
)

// interface to encapsulate any type of Metric
type Metric interface {
	String() string
	GetKey() string
	GetValue() interface{}
	SetKey(string)
	SetValue(interface{})
}

// Metric type to define the structure of a Counter
type Counter struct {
	key	string
	count	float64
}

// function to provide a string interface to the metric type
func (c *Counter) String() string {
	return fmt.Sprintf("%s\t%v",c.key,c.count)
}

// Function to return the key of the counter
func (c *Counter) GetKey() string {
	return c.key
}

// Function to return the value of the counter
func (c *Counter) GetValue() interface{} {
	return c.count
}

// function to set the key of a private variable
func (c *Counter) SetKey(key string) {
	c.key = key
}

// function to set and convert the private value variable
func (c *Counter) SetValue(value interface{}) {
	if val, ok := value.(float64); ok {
		c.count = val
	}
}

// Type to represent a time series logger
type TimeSeries struct {
	key	string
	time	time.Time
}

// pretty print the metric type
func (t *TimeSeries) String() string {
	return fmt.Sprintf("%s\t%v",t.key,t.time)
}

// Get the key for the given TimeSeries item
func (t *TimeSeries) GetKey() string {
	return t.key
}

// Get the value of the TimeSeries item
func (t *TimeSeries) GetValue() interface{} {
	return t.time
}

// set the key for the metric type
func (t *TimeSeries) SetKey(key string) {
	t.key = key
}

// set and convert the value for the TimeSeries value
func (t *TimeSeries) SetValue(value interface{}) {
	if val, ok := value.(time.Time); ok {
		t.time = val
	}
}

