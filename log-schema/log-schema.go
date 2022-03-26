package logSchema

import "github.com/diptamay/go-commons/glogger"

type LogSchema struct {
	glogger.DefaultLogMessage
	execRequestId string
	logiflowId    string
	miniApp       string
	clientId      string
	associateId   string
}
