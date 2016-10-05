package goroo

import (
	"bytes"
	"net"
	"strconv"
	"testing"
)

func BenchmarkToBytes(b *testing.B) {
	buffer := new(bytes.Buffer)
	gqtp := gqtpStruct{}
	gqtp.Protocol = 0xc7
	gqtp.QueryType = 0x02 // default is JSON
	gqtp.KeyLength = 0    // not used
	gqtp.Level = 0x00     // not used
	gqtp.Flags = cGrnGqtpFlagsTail
	gqtp.Status = 0 // not used
	gqtp.Size = 0
	gqtp.Opaque = 0 // not used
	gqtp.Cas = 0    // not used
	gqtp.Body = buffer.Bytes()

	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		gqtp.toByte(buffer)
	}
}

func BenchmarkClinetGqtpCall(b *testing.B) {
	body := byteBody
	mock := gqtpMock(body)
	defer mock.Close()

	host, p, _ := net.SplitHostPort(mock.Address)
	port, _ := strconv.Atoi(p)
	client := NewClient("gqtp", host, port)
	arg := map[string]string{}

	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.Call("table_list", arg)
		if err != nil {
			b.Fatal(err)
		}
	}
}
