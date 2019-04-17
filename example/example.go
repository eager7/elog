package main

import "github.com/eager7/elog"

func main() {
	log1 := elog.NewLogger("log1", elog.DebugLevel)
	log1.Debug("hello world!")
	log1.Notice("hello world!")
	log1.Info("hello world!")
	log1.Warn("hello world!")

	log2 := elog.NewLogger("log2", elog.WarnLevel)
	log2.Debug("hello world!")
	log2.Notice("hello world!")
	log2.Info("hello world!")
	log2.Warn("hello world!")
}
