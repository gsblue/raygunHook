package raygunHook

import (
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/gsblue/raygunclient/httpdata"
	"github.com/gsblue/raygunclient/stack"
)

const (
	//ErrorFieldName is the name of the field in logrus.Entry.Data, which should hold the error
	ErrorFieldName = "error"
	//RequestFieldName is the name of the field in logrus.Entry.Data, which should hold the request
	RequestFieldName = "request"
	//UserFieldName is the name of the field in logrus.Entry.Data, which should hold the user identifier
	UserFieldName = "user"
	//StackTraceFieldName is the name of the field in logrus.Entry.Data, which should hold stack trace
	StackTraceFieldName = "StackTrace"
	//CustomDataFieldName is the name of the field in logrus.Entry.Data, which should hold any custom data
	CustomDataFieldName = "customData"
)

//EntryWithRequest is a helper function to add request to a logrus.Entry
func EntryWithRequest(e *logrus.Entry, r *http.Request) *logrus.Entry {
	return e.WithField(RequestFieldName, httpdata.NewHTTPRequest(r))
}

//EntryWithUser is a helper function to add user identifier to a logrus Entry
func EntryWithUser(e *logrus.Entry, user string) *logrus.Entry {
	return e.WithField(UserFieldName, user)
}

//EntryWithCustomData is a helper function to add custom data to a logrus Entry
func EntryWithCustomData(e *logrus.Entry, data interface{}) *logrus.Entry {
	return e.WithField(CustomDataFieldName, data)
}

//EntryWithStackTrace is a helper function to add stack trace to a logrus Entry
func EntryWithStackTrace(e *logrus.Entry, trace stack.Trace) *logrus.Entry {
	return e.WithField(StackTraceFieldName, trace)
}
