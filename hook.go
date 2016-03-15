//Package raygunHook provides the hook for logrus, to send errors to raygun.io
//for more details visit https://www.github.com/gsblue/raygunHook
package raygunHook

import (
	"errors"

	"github.com/Sirupsen/logrus"
	ray "github.com/gsblue/raygunclient"
	"github.com/gsblue/raygunclient/httpdata"
	"github.com/gsblue/raygunclient/stack"
)

type hook struct {
	notifier    ray.Notifier
	defautltags []string
}

//HookConfig is the struct to hold the configuration values required for the hook.
type HookConfig struct {
	APIKey      string   //APIKey for your application in raygun. This field is mandatory.
	Version     string   //Version of your application
	DefaultTags []string //DefaultTags for your error entries
}

//NewHook creates a new raygun logrus.Hook
func NewHook(config *HookConfig) logrus.Hook {
	c := ray.NewClient(config.APIKey, config.Version, nil)

	return &hook{notifier: c, defautltags: config.DefaultTags}
}

//Fire sends the error from logrus.Entry to raygun
func (h *hook) Fire(e *logrus.Entry) error {
	var entry *ray.ErrorEntry

	if val, ok := e.Data[ErrorFieldName]; ok {
		if err, ok := val.(error); ok {
			entry = ray.NewErrorEntry(err)
		}
	}

	if entry == nil {
		entry = ray.NewErrorEntry(errors.New(e.Message))
	}

	if val, ok := e.Data[RequestFieldName]; ok {
		if req, ok := val.(*httpdata.HTTPRequest); ok {
			entry.Request = req
		}
	}

	if val, ok := e.Data[UserFieldName]; ok {
		if user, ok := val.(string); ok {
			entry.SetUser(user)
		}
	}

	if val, ok := e.Data[CustomDataFieldName]; ok {
		entry.SetCustomData(val)
	}

	entry.Tags = h.defautltags

	var trace stack.Trace
	if val, ok := e.Data[StackTraceFieldName]; ok {
		trace, _ = val.(stack.Trace)
	}

	if trace != nil {
		return h.notifier.NotifyWithStackTrace(entry, trace)
	}

	return h.notifier.Notify(entry)
}

//Levels returns the logrus.Level which this raygung hook supports
func (h *hook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}
