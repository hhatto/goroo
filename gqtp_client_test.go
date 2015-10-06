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
		t.Fatal(err)
	}
	if res.Status != 0 {
		t.Errorf("status not zero.[%d]", res.Status)
	}
	if len(res.Body.([]interface{})) != 1 {
		t.Errorf("body fail.[%s]", res.Body)
	}
}

func TestGqtp_TableList_Count1_Success(t *testing.T) {
	body := []byte{0xc7, 0x2, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x1, 0x17, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5b, 0x5b, 0x5b, 0x22, 0x69, 0x64, 0x22, 0x2c, 0x22, 0x55, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x70, 0x61, 0x74, 0x68, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x72, 0x61, 0x6e, 0x67, 0x65, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x69, 0x7a, 0x65, 0x72, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x2c, 0x5b, 0x22, 0x6e, 0x6f, 0x72, 0x6d, 0x61, 0x6c, 0x69, 0x7a, 0x65, 0x72, 0x22, 0x2c, 0x22, 0x53, 0x68, 0x6f, 0x72, 0x74, 0x54, 0x65, 0x78, 0x74, 0x22, 0x5d, 0x5d, 0x2c, 0x5b, 0x32, 0x35, 0x36, 0x2c, 0x22, 0x47, 0x51, 0x54, 0x50, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x22, 0x2c, 0x22, 0x2e, 0x2f, 0x6d, 0x61, 0x72, 0x6b, 0x64, 0x6f, 0x77, 0x6e, 0x2e, 0x64, 0x62, 0x2e, 0x30, 0x30, 0x30, 0x30, 0x31, 0x30, 0x30, 0x22, 0x2c, 0x22, 0x54, 0x41, 0x42, 0x4c, 0x45, 0x5f, 0x48, 0x41, 0x53, 0x48, 0x5f, 0x4b, 0x45, 0x59, 0x7c, 0x50, 0x45, 0x52, 0x53, 0x49, 0x53, 0x54, 0x45, 0x4e, 0x54, 0x22, 0x2c, 0x6e, 0x75, 0x6c, 0x6c, 0x2c, 0x6e, 0x75, 0x6c, 0x6c, 0x2c, 0x6e, 0x75, 0x6c, 0x6c, 0x2c, 0x6e, 0x75, 0x6c, 0x6c, 0x5d, 0x5d}
	mock := gqtpMock(body)
	defer mock.Close()

	client := newGqtpClient(mock.Address)
	res, err := client.Call("table_list", map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != 0 {
		t.Errorf("status not zero.[%d]", res.Status)
	}
	if len(res.Body.([]interface{})) != 2 {
		t.Errorf("body fail.[%s]", res.Body)
	}
}

func TestGqtp_ColumnCreate_UserName_Success(t *testing.T) {
	body := []byte{0xc7, 0x2, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x74, 0x72, 0x75, 0x65}
	mock := gqtpMock(body)
	defer mock.Close()

	client := newGqtpClient(mock.Address)
	res, err := client.Call("column_create", map[string]string{
		"table": "GQTPTable",
		"name":  "user_name",
		"type":  "ShortText",
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != 0 {
		t.Errorf("status not zero.[%d]", res.Status)
	}
	if res.Body.(bool) != true {
		t.Errorf("body fail.[%s]", res.Body)
	}
}

func TestGqtp_ColumnCreate_UserName_Fail(t *testing.T) {
	body := []byte{0xc7, 0x2, 0x0, 0x0, 0x0, 0x2, 0xff, 0xea, 0x0, 0x0, 0x0, 0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x66, 0x61, 0x6c, 0x73, 0x65}
	mock := gqtpMock(body)
	defer mock.Close()

	client := newGqtpClient(mock.Address)
	res, err := client.Call("column_create", map[string]string{
		"table": "GQTPTable",
		"name":  "user_name",
		"type":  "ShortText",
	})
	if err == nil {
		t.Errorf("err is nil")
	}
	if res != nil {
		t.Errorf("res is not nil")
	}
}
