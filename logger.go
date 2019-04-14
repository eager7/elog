package elog

import (
	"fmt"
	"github.com/BlockABC/wallet_eth_client/common/elog/bunnystub"
	"github.com/BlockABC/wallet_eth_client/common/elog/logbunny"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
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
	httpPort           string
	rollingTimePattern string
	skip               int
	logger             logbunny.Logger
}

var Log Logger //default logger
var levelHandler *logbunny.HTTPHandler

func Initialize(dir string) error {
	//fmt.Println("init log, baseDir:", dir)
	logOpt, err := newLoggerOpt()
	if err != nil {
		return err
	}

	logFilename := map[logbunny.LogLevel]string{
		logbunny.DebugLevel: dir + logOpt.debugLogFilename,
		logbunny.InfoLevel:  dir + logOpt.infoLogFilename,
		logbunny.WarnLevel:  dir + logOpt.warnLogFilename,
		logbunny.ErrorLevel: dir + logOpt.errorLogFilename,
		logbunny.FatalLevel: dir + logOpt.fatalLogFilename,
	}

	outputWriter := make(map[logbunny.LogLevel][]io.Writer)
	for level, file := range logFilename {
		logFileWriter, err := newLogFile(file, logOpt.rollingTimePattern)
		if err != nil {
			return err
		}
		if logOpt.toTerminal {
			outputWriter[level] = []io.Writer{logFileWriter, os.Stdout}
		} else {
			outputWriter[level] = []io.Writer{logFileWriter}
		}
	}

	zapCfg := &logbunny.Config{
		Type:        logOpt.loggerType,
		Level:       logOpt.debugLevel,
		Encoder:     logOpt.loggerEncoder,
		WithCaller:  logOpt.withCaller,
		Out:         nil,
		WithNoLock:  false,
		TimePattern: logOpt.timePattern,
		Skip:        logOpt.skip,
	}
	logOpt.logger, err = logbunny.FilterLogger(zapCfg, outputWriter)
	if err != nil {
		return err
	}

	logbunny.SetCallerSkip(3)
	// log.Warp()

	levelHandler = logbunny.NewHTTPHandler(logOpt.logger)
	http.HandleFunc("/logoutLevel", func(w http.ResponseWriter, r *http.Request) {
		levelHandler.ServeHTTP(w, r)
	})
	go func() {
		if err := http.ListenAndServe(logOpt.httpPort, nil); err != nil {
			fmt.Println(err)
		}
	}()
	Log = logOpt
	return nil
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
		httpPort:           viper.GetString("log.http_port"),
		rollingTimePattern: viper.GetString("log.rolling_time_pattern"),
		skip:               viper.GetInt("log.skip"),
		logger:             nil,
	}, nil
}

func InitDefaultConfig() {
	viper.SetDefault("log.debug_level", 0) //0: debug 1: info 2: warn 3: error 4: panic 5: fatal
	viper.SetDefault("log.loggerType", 0)  //0: zap 1: logrus
	viper.SetDefault("log.with_caller", false)
	viper.SetDefault("log.logger_encoder", 1) //0: json 1: console
	viper.SetDefault("log.time_pattern", "2006-01-02 15:04:05.00000")
	viper.SetDefault("log.http_port", ":50015")                 // RESTFul API to change logout level dynamically
	viper.SetDefault("log.debug_log_filename", "debug.log")       //or 'stdout' / 'stderr'
	viper.SetDefault("log.info_log_filename", "debug.log")        //or 'stdout' / 'stderr'
	viper.SetDefault("log.warn_log_filename", "warn.log")       //or 'stdout' / 'stderr'
	viper.SetDefault("log.error_log_filename", "err.log")       //or 'stdout' / 'stderr'
	viper.SetDefault("log.fatal_log_filename", "fatal.log")       //or 'stdout' / 'stderr'
	viper.SetDefault("log.rolling_time_pattern", "0 0 0 * * *") //rolling the log everyday at 00:00:00
	viper.SetDefault("log.skip", 4)                             //call depth, zap log is 3, logger is 4
	viper.SetDefault("log.to_terminal", true)                   //out put to terminal
}

func newLogFile(logPath string, rollingTimePattern string) (io.Writer, error) {
	if file := stdOutput(logPath); file != nil {
		return file, nil
	} else {
		file, err := bunnystub.NewIOWriter(logPath, bunnystub.TimeRotate, bunnystub.WithTimePattern(rollingTimePattern))
		if err != nil {
			return nil, err
		}
		return file, nil
	}
}
func stdOutput(logPath string) *os.File {
	if logPath == "stdout" {
		return os.Stdout
	}
	if logPath == "stderr" {
		return os.Stderr
	}
	return nil
}

const (
	colorRed = iota + 91
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
)

type Logger interface {
	Notice(a ...interface{})
	Debug(a ...interface{})
	Info(a ...interface{})
	Warn(a ...interface{})
	Error(a ...interface{})
	Fatal(a ...interface{})
	Panic(a ...interface{})
	GetLogger() logbunny.Logger
	ErrStack()
	SetLogLevel(level int) error
	//GetLogLevel() int
}

func (l *loggerOpt) Notice(a ...interface{}) {
	if l.loggerEncoder == 0 {
		l.logger.Debug(fmt.Sprintln(a...))
	} else {
		msg := "\x1b[" + strconv.Itoa(colorGreen) + "m" + "▶ " + fmt.Sprintln(a...) + "\x1b[0m"
		l.logger.Debug(msg)
	}
}

func (l *loggerOpt) Debug(a ...interface{}) {
	if l.loggerEncoder == 0 {
		l.logger.Debug(fmt.Sprintln(a...))
	} else {
		msg := "\x1b[" + strconv.Itoa(colorBlue) + "m" + "▶ " + fmt.Sprintln(a...) + "\x1b[0m"
		l.logger.Debug(msg)
	}
}

func (l *loggerOpt) Info(a ...interface{}) {
	if l.loggerEncoder == 0 {
		l.logger.Info(fmt.Sprintln(a...))
	} else {
		msg := "\x1b[" + strconv.Itoa(colorYellow) + "m" + "▶ " + fmt.Sprintln(a...) + "\x1b[0m"
		l.logger.Info(msg)
	}
}

func (l *loggerOpt) Warn(a ...interface{}) {
	if l.loggerEncoder == 0 {
		l.logger.Warn(fmt.Sprintln(a...))
	} else {
		msg := "\x1b[" + strconv.Itoa(colorMagenta) + "m" + "▶ " + fmt.Sprintln(a...) + "\x1b[0m"
		l.logger.Warn(msg)
	}
}

func (l *loggerOpt) Error(a ...interface{}) {
	if l.loggerEncoder == 0 {
		l.logger.Error(fmt.Sprintln(a...))
	} else {
		msg := "\x1b[" + strconv.Itoa(colorRed) + "m" + "▶ " + fmt.Sprintln(a...) + "\x1b[0m"
		l.logger.Error(msg)
	}
}

func (l *loggerOpt) Fatal(a ...interface{}) {
	if l.loggerEncoder == 0 {
		l.logger.Fatal(fmt.Sprintln(a...))
	} else {
		msg := "\x1b[" + strconv.Itoa(colorYellow) + "m" + "▶ " + fmt.Sprintln(a...) + "\x1b[0m"
		l.logger.Fatal(msg)
	}
}

func (l *loggerOpt) Panic(a ...interface{}) {
	if l.loggerEncoder == 0 {
		l.logger.Panic(fmt.Sprintln(a...))
	} else {
		msg := "\x1b[" + strconv.Itoa(colorYellow) + "m" + "▶ " + fmt.Sprintln(a...) + "\x1b[0m"
		l.logger.Panic(msg)
	}
	panic(fmt.Sprintln(a...))
}

func (l *loggerOpt) ErrStack() {
	l.Warn(string(debug.Stack()))
}

func (l *loggerOpt) SetLogLevel(level int) error {
	l.logger.SetLevel(logbunny.LogLevel(level))
	return nil
}

func (l *loggerOpt) GetLogger() logbunny.Logger {
	return l.logger
}
