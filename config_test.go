package goblog

import (
	"fmt"
	"testing"
)

func TestParseConf(t *testing.T) {
	conf, err := ParseConf("config.json")
	if err != nil {
		t.Error(err)
	}

	fmt.Print(conf)
}

func TestWriteConf(t *testing.T) {
	conf, err := ParseConf("config.json")
	if err != nil {
		t.Error(err)
	}

	err = WriteConf(conf, "config.json")

	if err != nil {
		t.Error(err)
	}
}
