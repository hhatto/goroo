package goroo

import (
	"fmt"
)

type GroongaResult struct {
	RawData     string
	Status      int
	StartTime   float64
	ElapsedTime float64
	Body        any
}

// Deprecated: It is scheduled to be abolished.
type GroongaClient struct {
	Protocol string
	Host     string
	Port     int
	client   Client
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

// Deprecated: This function is scheduled to be replaced with NewClient.
func NewGroongaClient(protocol, host string, port int) *GroongaClient {
	c := &GroongaClient{
		Protocol: protocol,
		Host:     host,
		Port:     port,
	}
	c.client = NewClient(protocol, host, port)
	return c
}

// Deprecated: This function is scheduled to be replaced with NewClient.
func (g *GroongaClient) Call(command string, params map[string]string) (*GroongaResult, error) {
	return g.client.Call(command, params)
}
