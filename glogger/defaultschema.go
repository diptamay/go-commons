package glogger

import "time"

type DefaultLogMessage struct {
	action      string
	containerId string
	duration    int
	filename    string
	fromSpanId  string
	kind        string
	level       string
	method      string
	msg         string
	name        string
	path        string
	pid         int
	schema      string
	severity    int
	spanId      string
	startTime   int64
	status      string
	time        time.Time
	traceId     string
	v           string
	version     string
}
