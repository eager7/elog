## elog
从uber改造的log库

## 获取
go get

##使用
new一个log出来用即可，会在程序执行目录下生成一个elog.toml的配置文件，自行修改对应参数即可
```
func TestLogger(t *testing.T) {
	l := elog.NewLogger("example", 0)
	l.Debug("debug ------------------")
	l.Info("info   --------------------")
	l.Warn("warn   ----------------")
	l.Error("error ---------------------")
	//l.Panic("panic------------")
	//defaultLog.Log.ErrStack()
}
```
##参数
下面是默认参数
```
	viper.SetDefault("log.debug_level", 0) //0: debug 1: info 2: warn 3: error 4: panic 5: fatal
	viper.SetDefault("log.loggerType", 0)  //0: zap 1: logrus
	viper.SetDefault("log.with_caller", false)
	viper.SetDefault("log.logger_encoder", 1) //0: json 1: console
	viper.SetDefault("log.time_pattern", "2006-01-02 15:04:05.00000")
	viper.SetDefault("log.http_port", ":50015")                 // RESTFul API to change logout level dynamically
	viper.SetDefault("log.debug_log_filename", "debug.log")     //or 'stdout' / 'stderr'
	viper.SetDefault("log.info_log_filename", "info.log")       //or 'stdout' / 'stderr'
	viper.SetDefault("log.warn_log_filename", "warn.log")       //or 'stdout' / 'stderr'
	viper.SetDefault("log.error_log_filename", "err.log")       //or 'stdout' / 'stderr'
	viper.SetDefault("log.fatal_log_filename", "fatal.log")     //or 'stdout' / 'stderr'
	viper.SetDefault("log.rolling_time_pattern", "0 0 0 * * *") //rolling the log everyday at 00:00:00
	viper.SetDefault("log.skip", 4)                             //call depth, zap log is 3, logger is 4
	viper.SetDefault("log.to_terminal", true)                   //out put to terminal
  ```
