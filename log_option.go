package elog

import (
	"github.com/eager7/elog/logbunny"
	"github.com/spf13/viper"
)

type loggerOpt struct {
	debugLevel         logbunny.LogLevel
	loggerType         logbunny.LogType
	withCaller         bool
	toTerminal         bool
	loggerEncoder      logbunny.EncoderType
	timePattern        string
	debugLogFilename   string
	infoLogFilename    string
	warnLogFilename    string
	errorLogFilename   string
	fatalLogFilename   string
	rollingTimePattern string
	skip               int
	logger             logbunny.Logger
}

func newLoggerOpt() (*loggerOpt, error) {
	return &loggerOpt{
		debugLevel:         logbunny.LogLevel(viper.GetInt("log.debug_level")),
		loggerType:         logbunny.LogType(viper.GetInt("log.loggerType")),
		withCaller:         viper.GetBool("log.with_caller"),
		toTerminal:         viper.GetBool("log.to_terminal"),
		loggerEncoder:      logbunny.EncoderType(viper.GetInt("log.logger_encoder")),
		timePattern:        viper.GetString("log.time_pattern"),
		debugLogFilename:   viper.GetString("log.debug_log_filename"),
		infoLogFilename:    viper.GetString("log.info_log_filename"),
		warnLogFilename:    viper.GetString("log.warn_log_filename"),
		errorLogFilename:   viper.GetString("log.error_log_filename"),
		fatalLogFilename:   viper.GetString("log.fatal_log_filename"),
		rollingTimePattern: viper.GetString("log.rolling_time_pattern"),
		skip:               viper.GetInt("log.skip"),
		logger:             nil,
	}, nil
}

func init() {
	viper.SetDefault("log.debug_level", 0) //0: debug 1: info 2: warn 3: error 4: panic 5: fatal
	viper.SetDefault("log.loggerType", 0)  //0: zap 1: logrus
	viper.SetDefault("log.with_caller", false)
	viper.SetDefault("log.logger_encoder", 1) //0: json 1: console
	viper.SetDefault("log.time_pattern", "2006-01-02 15:04:05.00000")
	viper.SetDefault("log.http_port", ":50015")                 // RESTFul API to change logout level dynamically
	viper.SetDefault("log.debug_log_filename", "debug.log")     //or 'stdout' / 'stderr'
	viper.SetDefault("log.info_log_filename", "debug.log")      //or 'stdout' / 'stderr'
	viper.SetDefault("log.warn_log_filename", "warn.log")       //or 'stdout' / 'stderr'
	viper.SetDefault("log.error_log_filename", "err.log")       //or 'stdout' / 'stderr'
	viper.SetDefault("log.fatal_log_filename", "fatal.log")     //or 'stdout' / 'stderr'
	viper.SetDefault("log.rolling_time_pattern", "0 0 0 * * *") //rolling the log everyday at 00:00:00
	viper.SetDefault("log.skip", 4)                             //call depth, zap log is 3, logger is 4
	viper.SetDefault("log.to_terminal", true)                   //out put to terminal
}
