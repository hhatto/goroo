package goroo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
)

type gqtpClient struct {
	address string
}

func (c *gqtpClient) Call(command string, params map[string]string) (*GroongaResult, error) {
	body, err := callGQTP(c.address, command, params)
	if err != nil {
		return nil, err
	}
	res, err := setResult(body)
	return &res, nil
}

func callGQTP(address, command string, params map[string]string) (b []byte, err error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return b, err
	}
	defer conn.Close()

	buffer := bytes.NewBufferString(command)
	for value, name := range params {
		buffer.WriteString(fmt.Sprintf(" --%s '%s'", value, name))
	}
	bodyLen := uint32(len(buffer.String()))

	// encode request header and body
	gqtp := GQTP{}
	gqtp.Protocol = 0xc7
	gqtp.QueryType = 0x02            // default is JSON
	gqtp.KeyLength = make([]byte, 2) // not used
	gqtp.Level = 0x00                // not used
	gqtp.Flags = GRN_GQTP_FLAGS_TAIL
	gqtp.Status = make([]byte, 2) // not used
	gqtp.Size = []byte{
		byte(0xff000000 & bodyLen),
		byte(0x00ff0000 & bodyLen),
		byte(0x0000ff00 & bodyLen),
		byte(0x000000ff & bodyLen),
	}
	gqtp.Opaque = make([]byte, 4) // not used
	gqtp.Cas = make([]byte, 8)    // not used
	gqtp.Body = buffer.Bytes()

	_, err = conn.Write(gqtp.toByte())
	if err != nil {
		return b, err
	}

	// TODO: recieve over 1024 byte
	resp := make([]byte, 1024)
	nr, err := conn.Read(resp)
	if err != nil {
		return b, err
	}

	// decode respose header and body
	if 0xc7 != byte(resp[0]) {
		return b, fmt.Errorf("check response protocol NG 0x%x", resp[0])
	}
	if GRN_GQTP_FLAGS_TAIL != resp[5] {
		return b, fmt.Errorf("flag: %v is not support", resp[5])
	}
	status := uint32(resp[7]) + uint32(resp[6])<<8
	if status != 0 {
		return b, fmt.Errorf("status error: [%d]%s", status, GRN_GQTP_STATUS[status])
	}
	respBodyLen := (uint32(resp[8])<<24)&0xff000000 +
		(uint32(resp[9])<<16)&0x00ff0000 +
		(uint32(resp[10])<<8)&0x0000ff00 +
		(uint32(resp[11]))&0x000000ff
	if int(respBodyLen) != nr-GRN_GQTP_HEADER_SIZE {
		return b, fmt.Errorf("invalid body size: [%d]", respBodyLen)
	}

	return resp[GRN_GQTP_HEADER_SIZE:nr], err
}

func setResult(body []byte) (result GroongaResult, err error) {
	result.RawData = string(body)

	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return result, err
	}

	result.Body = data
	return result, nil
}

func newGqtpClient(address string) Client {
	return &gqtpClient{address}
}
