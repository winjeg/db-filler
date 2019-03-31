package log

import (
	"github.com/sirupsen/logrus"
	"github.com/winjeg/db-filler/config"
	"github.com/winjeg/go-commons/log"
)

var logger = log.GetLogger(config.GetConf())


// get logger
func GetLogger() *logrus.Logger {
	return logger
}
