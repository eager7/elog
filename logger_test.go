package elog_test

import (
	"github.com/eager7/elog"
	"testing"
)

func TestLogger(t *testing.T) {
	l := elog.NewLogger("example", 0)
	l.Debug("debug ------------------")
	l.Info("info   --------------------")
	l.Warn("warn   ----------------")
	l.Error("error ---------------------")
	//l.Panic("panic------------")
	//defaultLog.Log.ErrStack()
}
