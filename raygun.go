package raygun

import (
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/sditools/goraygun"
)

type raygunHook struct {
	Client *goraygun.Client
}

func NewHook(endpoint string, apiKey string, env string) *raygunHook {
	client := goraygun.Init(goraygun.Settings{
		ApiKey:      apiKey,
		Endpoint:    endpoint,
		Environment: env,
	}, goraygun.Entry{})
	return &raygunHook{client}
}

func (hook *raygunHook) Fire(logEntry *logrus.Entry) error {
	// Start with a copy of the default entry
	raygunEntry := hook.Client.Entry

	if request, ok := logEntry.Data["request"]; ok {
		raygunEntry.Details.Request.Populate(*(request.(*http.Request)))
	}

	var reportErr error
	if err, ok := logEntry.Data["error"]; ok {
		reportErr = err.(error)
	} else {
		reportErr = errors.New(logEntry.Message)
	}

	hook.Client.Report(reportErr, raygunEntry)

	return nil
}

func (hook *raygunHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}
