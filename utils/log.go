package utils

import (
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"path"
	"time"
	log "github.com/sirupsen/logrus"
)

func GetLogger(fileName string)*log.Logger {
	logger := log.New()
	ConfigLocalFilesystemLogger(logger,"/Users/king/go/framework/bin/log/"+fileName, fileName, time.Hour*24*60)
	return logger
}

func LogMessage(useLog *log.Logger, msg string, mcmd uint16, scmd uint16, data []byte) {
	useLog.WithFields(log.Fields{
		"mcmd": mcmd,
		"scmd": scmd,
		"data": string(data),
	}).Info(msg)
}

func ConfigLocalFilesystemLogger(logger *log.Logger,logPath string, logFileName string, maxAge time.Duration) {
	baseLogPaht := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPaht+".%Y%m%d.log",
		rotatelogs.WithLinkName(baseLogPaht),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Hour*24), // 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, &log.JSONFormatter{
		TimestampFormat:"2006-01-02 15:04:05",
	})
	logger.AddHook(lfHook)
}