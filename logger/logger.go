package logger

import (
	"time"

	"github.com/sirupsen/logrus"
	formatter "github.com/tim-ywliu/nested-logrus-formatter"
)

var (
	log         *logrus.Logger
	NasLog      *logrus.Entry
	NasMsgLog   *logrus.Entry
	ConvertLog  *logrus.Entry
	SecurityLog *logrus.Entry
)

func init() {
	log = logrus.New()
	log.SetReportCaller(false)

	log.Formatter = &formatter.Formatter{
		TimestampFormat: time.RFC3339,
		TrimMessages:    true,
		NoFieldsSpace:   true,
		HideKeys:        true,
		FieldsOrder:     []string{"component", "category"},
	}

	NasLog = log.WithFields(logrus.Fields{"component": "LIB", "category": "NAS"})
	NasMsgLog = log.WithFields(logrus.Fields{"component": "NAS", "category": "Message"})
	ConvertLog = log.WithFields(logrus.Fields{"component": "NAS", "category": "Convert"})
	SecurityLog = log.WithFields(logrus.Fields{"component": "NAS", "category": "Security"})
}

func GetLogger() *logrus.Logger {
	return log
}

func SetLogLevel(level logrus.Level) {
	NasLog.Infoln("set log level :", level)
	log.SetLevel(level)
}

func SetReportCaller(enable bool) {
	NasLog.Infoln("set report call :", enable)
	log.SetReportCaller(enable)
}
