package raygunHook

import (
	"github.com/gsblue/raygunHook/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	ray "github.com/gsblue/raygunHook/Godeps/_workspace/src/github.com/gsblue/raygun4go"
	"net/http"
)

type hook struct {
	Client raygunClient
}

//HookConfig is the struct to hold the configuration values required for the hook.
type HookConfig struct {
	APIKey  string   //ApiKey for your raygun account. This field is mandatory.
	AppName string   //AppName is your application name. This field is mandatory.
	Version string   //Version of your application
	Tags    []string //Tags which get added to all the error entries
}

type raygunClient interface {
	CreateErrorEntry(err error) *ray.ErrorEntry
	CreateErrorEntryFromMsg(msg string) *ray.ErrorEntry
	SubmitError(entry *ray.ErrorEntry) error
}

const (
	//ErrorFieldName is the name of the field in logrus.Entry.Data, which should hold the error
	ErrorFieldName = "error"
	//RequestFieldName is the name of the field in logrus.Entry.Data, which should hold the request
	RequestFieldName = "request"
	//UserFieldName is the name of the field in logrus.Entry.Data, which should hold the user identifier
	UserFieldName = "user"
	//CustomDataFieldName is the name of the field in logrus.Entry.Data, which should hold any custom data
	CustomDataFieldName = "customData"
)

//NewHook creates a new raygun logrus.Hook
func NewHook(config *HookConfig) (logrus.Hook, error) {
	c, err := ray.New(config.AppName, config.APIKey)
	if err != nil {
		return nil, err
	}
	c.Version(config.Version).Tags(config.Tags)

	return &hook{Client: c}, nil
}

//Fire sends the error from logrus.Entry to raygun
func (h *hook) Fire(e *logrus.Entry) error {
	var entry *ray.ErrorEntry

	if val, ok := e.Data[ErrorFieldName]; ok {
		if err, ok := val.(error); ok {
			entry = h.Client.CreateErrorEntry(err)
		}
	}

	if entry == nil {
		entry = h.Client.CreateErrorEntryFromMsg(e.Message)
	}

	if val, ok := e.Data[RequestFieldName]; ok {
		if req, ok := val.(*http.Request); ok {
			entry.SetRequest(req)
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

	return h.Client.SubmitError(entry)
}

//Levels returns the logrus.Level which this raygung hook supports
func (h *hook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}
