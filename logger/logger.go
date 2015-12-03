/*
Package logger init and export 'Log'.

根据环境变量 GO_RABBIT_ENV 来设置 Log.Formatter.

"production": Formatter 为 'JSONFormatter'

其它: Formatter 为 'TextFormatter'
*/
package logger

import (
	"github.com/Sirupsen/logrus"
	"os"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	if logrusLevel := os.Getenv("RABBIT_GO_ENV"); logrusLevel == "production" {
		Log.Formatter = &logrus.JSONFormatter{}
	} else {
		Log.Formatter = &logrus.TextFormatter{}
	}
}
