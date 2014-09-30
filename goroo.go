package goroo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
)

const (
	GRN_GQTP_FLAGS_MORE  byte = 0x01
	GRN_GQTP_FLAGS_TAIL       = 0x02
	GRN_GQTP_FLAGS_HEAD       = 0x04
	GRN_GQTP_FLAGS_QUIET      = 0x08
	GRN_GQTP_FLAGS_QUIT       = 0x10
)
const GRN_GQTP_HEADER_SIZE int = 24

var GRN_GQTP_STATUS map[uint32]string = map[uint32]string{
	0:     "SUCCESS",
	1:     "END_OF_DATA",
	65535: "UNKNOWN_ERROR",
	65534: "OPERATION_NOT_PERMITTED",
	65533: "NO_SUCH_FILE_OR_DIRECTORY",
	65532: "NO_SUCH_PROCESS",
	65531: "INTERRUPTED_FUNCTION_CALL",
	65530: "INPUT_OUTPUT_ERROR",
	65529: "NO_SUCH_DEVICE_OR_ADDRESS",
	65528: "ARG_LIST_TOO_LONG",
	65527: "EXEC_FORMAT_ERROR",
	65526: "BAD_FILE_DESCRIPTOR",
	65525: "NO_CHILD_PROCESSES",
	65524: "RESOURCE_TEMPORARILY_UNAVAILABLE",
	65523: "NOT_ENOUGH_SPACE",
	65522: "PERMISSION_DENIED",
	65521: "BAD_ADDRESS",
	65520: "RESOURCE_BUSY",
	65519: "FILE_EXISTS",
	65518: "IMPROPER_LINK",
	65517: "NO_SUCH_DEVICE",
	65516: "NOT_A_DIRECTORY",
	65515: "IS_A_DIRECTORY",
	65514: "INVALID_ARGUMENT",
	65513: "TOO_MANY_OPEN_FILES_IN_SYSTEM",
	65512: "TOO_MANY_OPEN_FILES",
	65511: "INAPPROPRIATE_I_O_CONTROL_OPERATION",
	65510: "FILE_TOO_LARGE",
	65509: "NO_SPACE_LEFT_ON_DEVICE",
	65508: "INVALID_SEEK",
	65507: "READ_ONLY_FILE_SYSTEM",
	65506: "TOO_MANY_LINKS",
	65505: "BROKEN_PIPE",
	65504: "DOMAIN_ERROR",
	65503: "RESULT_TOO_LARGE",
	65502: "RESOURCE_DEADLOCK_AVOIDED",
	65501: "NO_MEMORY_AVAILABLE",
	65500: "FILENAME_TOO_LONG",
	65499: "NO_LOCKS_AVAILABLE",
	65498: "FUNCTION_NOT_IMPLEMENTED",
	65497: "DIRECTORY_NOT_EMPTY",
	65496: "ILLEGAL_BYTE_SEQUENCE",
	65495: "SOCKET_NOT_INITIALIZED",
	65494: "OPERATION_WOULD_BLOCK",
	65493: "ADDRESS_IS_NOT_AVAILABLE",
	65492: "NETWORK_IS_DOWN",
	65491: "NO_BUFFER",
	65490: "SOCKET_IS_ALREADY_CONNECTED",
	65489: "SOCKET_IS_NOT_CONNECTED",
	65488: "SOCKET_IS_ALREADY_SHUTDOWNED",
	65487: "OPERATION_TIMEOUT",
	65486: "CONNECTION_REFUSED",
	65485: "RANGE_ERROR",
	65484: "TOKENIZER_ERROR",
	65483: "FILE_CORRUPT",
	65482: "INVALID_FORMAT",
	65481: "OBJECT_CORRUPT",
	65480: "TOO_MANY_SYMBOLIC_LINKS",
	65479: "NOT_SOCKET",
	65478: "OPERATION_NOT_SUPPORTED",
	65477: "ADDRESS_IS_IN_USE",
	65476: "ZLIB_ERROR",
	65475: "LZO_ERROR",
	65474: "STACK_OVER_FLOW",
	65473: "SYNTAX_ERROR",
	65472: "RETRY_MAX",
	65471: "INCOMPATIBLE_FILE_FORMAT",
	65470: "UPDATE_NOT_ALLOWED",
	65469: "TOO_SMALL_OFFSET",
	65468: "TOO_LARGE_OFFSET",
	65467: "TOO_SMALL_LIMIT",
	65466: "CAS_ERROR",
	65465: "UNSUPPORTED_COMMAND_VERSION",
}

type GQTP struct {
	Protocol  byte
	QueryType byte
	KeyLength []byte // 2byte
	Level     byte
	Flags     byte
	Status    []byte // 2byte
	Size      []byte // 4byte
	Opaque    []byte // 4byte
	Cas       []byte // 8byte
	Body      []byte
}

func (gqtp *GQTP) toByte() (b []byte) {
	buffer := bytes.NewBuffer(b)
	buffer.WriteByte(gqtp.Protocol)
	buffer.WriteByte(gqtp.QueryType)
	buffer.Write(gqtp.KeyLength)
	buffer.WriteByte(gqtp.Level)
	buffer.WriteByte(gqtp.Flags)
	buffer.Write(gqtp.Status)
	buffer.Write(gqtp.Size)
	buffer.Write(gqtp.Opaque)
	buffer.Write(gqtp.Cas)
	buffer.Write(gqtp.Body)
	return buffer.Bytes()
}

var doGet = http.Get

type GroongaClient struct {
	Protocol string
	Host     string
	Port     int
}

type GroongaResult struct {
	RawData     string
	Status      int
	StartTime   float64
	ElapsedTime float64
	Body        interface{}
}

func NewGroongaClient(protocol, host string, port int) *GroongaClient {
	client := &GroongaClient{
		Protocol: protocol,
		Host:     host,
		Port:     port,
	}
	return client
}

func (client *GroongaClient) callGQTP(command string, params map[string]string) (b []byte, err error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.Host, client.Port))
	if err != nil {
		log.Println("Dial error:", err)
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
		log.Println("read error %v", err)
		return b, err
	}

	// decode respose header and body
	if 0xc7 != byte(resp[0]) {
		return b, fmt.Errorf("check response protocol NG 0x%x", resp[0])
	}
	//respGQTP.QueryType = byte(resp[1])
	//respGQTP.Flags = resp[5]
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

func (client *GroongaClient) callHTTP(command string, params map[string]string) (b []byte, err error) {
	v := url.Values{}
	for value, name := range params {
		v.Set(value, name)
	}
	requestUrl := fmt.Sprintf("%s://%s:%d/d/%s?%s",
		client.Protocol, client.Host, client.Port, command, v.Encode())
	resp, err := doGet(requestUrl)
	if err != nil {
		return nil, fmt.Errorf("http.Get() error: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response read error: %v", err)
	}

	return body, err
}

func (client *GroongaClient) Call(command string, params map[string]string) (result GroongaResult, err error) {
	if len(params) == 0 {
		return result, nil
	}

	var body []byte
	if client.Protocol == "gqtp" {
		// GQTP
		body, err = client.callGQTP(command, params)
		fmt.Println(string(body))
	} else {
		// HTTP
		body, err = client.callHTTP(command, params)
	}
	if err != nil {
		log.Println(err)
		return result, err
	}

	result, err = client.setResult(body)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (client *GroongaClient) setResult(body []byte) (result GroongaResult, err error) {
	result.RawData = fmt.Sprintf("%s", body)

	var data interface{}
	dec := json.NewDecoder(strings.NewReader(result.RawData))
	dec.Decode(&data)

	if client.Protocol == "http" {
		grnInfo := data.([]interface{})
		grnHeader := grnInfo[0].([]interface{})
		result.Status = int(grnHeader[0].(float64))
		result.StartTime = grnHeader[1].(float64)
		result.ElapsedTime = grnHeader[2].(float64)
		if len(grnHeader) == 3 {
			// groonga response ok
			result.Body = grnInfo[1]
		} else {
			// groonga response ng
			result.Body = grnHeader[3]
		}
	} else {
		result.Body = data.([]interface{})
	}

	return result, nil
}
