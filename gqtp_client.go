package goroo

type gqtpClient struct {
}

func (c *gqtpClient) Call(command string, params map[string]string) (*GroongaResult, error) {
	return &GroongaResult{}, nil
}

func newGqtpClient(con string) Client {
	return &gqtpClient{}
}
