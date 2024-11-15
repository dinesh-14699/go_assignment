package logger

import (
    "github.com/sirupsen/logrus"
)

var Log *logrus.Entry 

func InitLogger(url string, serviceType string) {
    baseLogger := logrus.New()

    httpHook := NewHTTPHook(url)
    baseLogger.AddHook(httpHook)

    baseLogger.SetFormatter(&logrus.JSONFormatter{})
    baseLogger.SetLevel(logrus.InfoLevel)

    Log = baseLogger.WithField("service_type", serviceType)
}

func UpdateLogContext(username string, userID string) {
    Log = Log.WithFields(logrus.Fields{
        "username": username,
        "user_id":  userID,
    })
}