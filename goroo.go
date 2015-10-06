package goroo

import (
	"fmt"
)

type GroongaResult struct {
	RawData     string
	Status      int
	StartTime   float64
	ElapsedTime float64
	Body        interface{}
}

type Client interface {
	Call(command string, params map[string]string) (*GroongaResult, error)
}

func NewClient(protocol, host string, port int) Client {
	if protocol == "http" {
		return newHttpClient(fmt.Sprintf("%s://%s:%d", protocol, host, port))
	}
	if protocol == "gqtp" {
		return newGqtpClient(fmt.Sprintf("%s:%d", host, port))
	}
	return nil
}
