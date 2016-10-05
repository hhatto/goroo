package goroo

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

const (
	cGrnGqtpFlagsMore  byte = 0x01
	cGrnGqtpFlagsTail       = 0x02
	cGrnGqtpFlagsHead       = 0x04
	cGrnGqtpFlagsQuiet      = 0x08
	cGrnGqtpFlagsQuit       = 0x10
)
const cGrnGqtpHeaderSize int = 24

var grnReturnCode map[uint32]string = map[uint32]string{
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

type gqtpStruct struct {
	Protocol  byte
	QueryType byte
	KeyLength uint16 // 2byte
	Level     byte
	Flags     byte
	Status    uint16 // 2byte
	Size      uint32 // 4byte
	Opaque    uint32 // 4byte
	Cas       uint64 // 8byte
	Body      []byte
}

func (gqtp *gqtpStruct) toByte(buffer *bytes.Buffer) []byte {
	scratch := make([]byte, cGrnGqtpHeaderSize)
	scratch[0] = gqtp.Protocol
	scratch[1] = gqtp.QueryType
	binary.BigEndian.PutUint16(scratch, gqtp.KeyLength)
	scratch[4] = gqtp.Level
	scratch[5] = gqtp.Flags
	binary.BigEndian.PutUint16(scratch, gqtp.Status)
	binary.BigEndian.PutUint32(scratch, gqtp.Size)
	binary.BigEndian.PutUint32(scratch, gqtp.Opaque)
	binary.BigEndian.PutUint64(scratch, gqtp.Cas)

	buffer.Write(scratch)
	buffer.Write(gqtp.Body)
	return buffer.Bytes()
}

type gqtpClient struct {
	address string
}

func (c *gqtpClient) Call(command string, params map[string]string) (*GroongaResult, error) {
	body, err := c.run(c.address, command, params)
	if err != nil {
		return nil, err
	}
	res, err := c.parse(body)
	return &res, nil
}

func (c *gqtpClient) run(address, command string, params map[string]string) (b []byte, err error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return b, err
	}
	defer conn.Close()

	buffer := bytes.NewBufferString(command)
	for value, name := range params {
		buffer.WriteString(fmt.Sprintf(" --%s '%s'", value, name))
	}

	// encode request header and body
	gqtp := gqtpStruct{}
	gqtp.Protocol = 0xc7
	gqtp.QueryType = 0x02 // default is JSON
	gqtp.KeyLength = 0    // not used
	gqtp.Level = 0x00     // not used
	gqtp.Flags = cGrnGqtpFlagsTail
	gqtp.Status = 0x00 // not used
	gqtp.Size = uint32(buffer.Len())
	gqtp.Opaque = 0 // not used
	gqtp.Cas = 0    // not used
	gqtp.Body = buffer.Bytes()

	buf := Buffs.Get().(*bytes.Buffer)
	msg := gqtp.toByte(buf)
	Buffs.Put(buf)

	_, err = conn.Write(msg)
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
	if cGrnGqtpFlagsTail != resp[5] {
		return b, fmt.Errorf("flag: %v is not support", resp[5])
	}
	status := uint32(resp[7]) + uint32(resp[6])<<8
	if status != 0 {
		return b, fmt.Errorf("status error: [%d]%s", status, grnReturnCode[status])
	}
	respBodyLen := (uint32(resp[8])<<24)&0xff000000 +
		(uint32(resp[9])<<16)&0x00ff0000 +
		(uint32(resp[10])<<8)&0x0000ff00 +
		(uint32(resp[11]))&0x000000ff
	if int(respBodyLen) != nr-cGrnGqtpHeaderSize {
		return b, fmt.Errorf("invalid body size: [%d]", respBodyLen)
	}

	return resp[cGrnGqtpHeaderSize:nr], err
}

func (c *gqtpClient) parse(body []byte) (result GroongaResult, err error) {
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
