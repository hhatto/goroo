package goroo

import (
	"fmt"
	"testing"
)

func TestGqtp_TableList_Empty_Success(t *testing.T) {
	client := newGqtpClient(fmt.Sprintf("%s:%d", "localhost", 10043))
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
