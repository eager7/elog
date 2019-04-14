package elog_test

import (
	"github.com/BlockABC/wallet_eth_client/common/elog"
	"testing"
)

func TestLogger(t *testing.T) {
	elog.InitDefaultConfig()
	if err := elog.Initialize("/tmp/"); err != nil {
		t.Fatal(err)
	}
	elog.Log.Debug("debug ------------------")
	elog.Log.Info("info   --------------------")
	elog.Log.Warn("warn   ----------------")
	elog.Log.Error("error ---------------------")
	//l.Panic("panic------------")
	//elog.Log.ErrStack()
}
