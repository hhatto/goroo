package goroo

import (
	"net"
	"testing"
)

func gqtpMock(body []byte) *gqtpServer {
	gs := newGqtpServer(func(conn net.Conn) {
		defer conn.Close()
		conn.Write(body)
	})
	return gs
}

func TestGqtp_TableList_Empty_Success(t *testing.T) {
	body := []byte{0xc7, 0x2, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0xbd, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5b, 0x5b, 0x5b, 0x22, 0x69, 0x64, 0x22, 0x2c, 0x22, 0x55, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x70, 0x61, 0x74, 0x68, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x69, 0x7a, 0x65, 0x72, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x6e, 0x6f, 0x72, 0x6d, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x72, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x5d, 0x5d}
	mock := gqtpMock(body)
	defer mock.Close()

	client := newGqtpClient(mock.Address)
	res, err := client.Call("table_list", map[string]string{})
	if err != nil {
		t.Error(err)
	}
	if res.Status != 0 {
		t.Errorf("status not zero.[%d]", res.Status)
	}
	if len(res.Body.([]interface{})) != 1 {
		t.Errorf("body fail.[%s]", res.Body)
	}
}
