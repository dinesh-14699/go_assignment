package logger

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type HTTPHook struct {
    URL string
}

func NewHTTPHook(url string) *HTTPHook {
    return &HTTPHook{URL: url}
}

func (hook *HTTPHook) Fire(entry *logrus.Entry) error {

    logData := make(map[string]interface{})

    logData["level"] = entry.Level.String()
    logData["message"] = entry.Message

    for key, value := range entry.Data {
        logData[key] = value
    }

    jsonData, err := json.Marshal(logData)
    if err != nil {
        return err
    }

    req, err := http.NewRequest("POST", hook.URL, bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }


    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    _, err = client.Do(req)
    return err
}

func (hook *HTTPHook) Levels() []logrus.Level {
    return logrus.AllLevels
}
