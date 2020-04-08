package core

import (
	"errors"
	"github.com/cloudflare/tableflip"
)

var Upg *tableflip.Upgrader

func UpgraderGoBlog() error {
	if Upg == nil {
		return errors.New("upg is nil")
	}
	return Upg.Upgrade()
}
