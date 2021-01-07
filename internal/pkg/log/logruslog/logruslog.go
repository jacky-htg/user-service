package logruslog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// UTCFormatter struct
type UTCFormatter struct {
	logrus.Formatter
}

// Format func
func (u UTCFormatter) Format(e *logrus.Entry) ([]byte, error) {
	e.Time = e.Time.UTC()
	return u.Formatter.Format(e)
}

// Init func
func Init() *logrus.Entry {

	log := logrus.New()

	logLocation := filepath.Join("./log/fluent-bit", os.Getenv("SERVICE_NAME"), fmt.Sprintf("%s.log", os.Getenv("SERVICE_NAME")))

	writer, _ := rotatelogs.New(
		logLocation+".%Y%m%d",
		rotatelogs.WithLinkName(logLocation),
		rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(604800)*time.Second),
	)

	log.Hooks.Add(lfshook.NewHook(
		lfshook.WriterMap{
			logrus.DebugLevel: writer,
		},
		&logrus.JSONFormatter{},
	))

	log.SetOutput(writer)

	logFile, err := os.OpenFile(logLocation, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file %s for output: %s", logLocation, err)
	}
	log.SetOutput(io.MultiWriter(os.Stderr, logFile))

	log.SetFormatter(UTCFormatter{&logrus.JSONFormatter{}})

	if os.Getenv("ENV") == "production" {
		log.SetLevel(logrus.InfoLevel)
	} else {
		log.SetLevel(logrus.DebugLevel)
	}

	return log.WithFields(
		logrus.Fields{
			"service": fmt.Sprintf("%s-%s", os.Getenv("SERVICE_NAME"), os.Getenv("SERVICE_VERSION")),
		})
}
