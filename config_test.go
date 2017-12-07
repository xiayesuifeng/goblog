package goblog

import (
	"testing"
	"fmt"
)

func TestParseConf(t *testing.T) {
	conf,err := ParseConf("config.json")
	if err!= nil {
		t.Error(err)
	}

	fmt.Print(conf)
}