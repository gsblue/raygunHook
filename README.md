# Raygun Hook [![Build Status](https://travis-ci.org/gsblue/raygunHook.svg)](https://travis-ci.org/gsblue/raygunHook) [![Coverage Status](https://coveralls.io/repos/github/gsblue/raygunHook/badge.svg?branch=master)](https://coveralls.io/github/gsblue/raygunHook?branch=master)

A Raygun.io hook for logrus. This package uses fork of [raygun4go](https://github.com/gsblue/raygun4go) client which allows concurrent calls to submit errors to raygun.

##Usage

```go
import (
	log "github.com/Sirupsen/logrus"
	"github.com/gsblue/raygunHook"
)

func init() {
	h, err := NewHook(&HookConfig{
		ApiKey:  "0Ih1ot2kApNL26wK0OAdvA==",
		AppName: "raygun hook test",
		Version: "2.1.10",
		Tags:    []string{"development"},
	})

	if err != nil {
		panic(err)
	}

	log.AddHook(h)
}

func SomeFunctionWhichLogs() {
	
	err := errors.New("some error")
	r, _ := http.NewRequest("GET", "http://www.google.com", nil)
	log.WithError(err).
		WithField(RequestFieldName, r). //to ensure request is sent to raygun
		WithField(UserFieldName, "john doe"). //to ensure user identifier is sent to raygun
		WithField(CustomDataFieldName, &struct{ OrderNo int }{340}). //to ensure custom data is sent to raygun
		Error()
}
```

Pull requests are welcome.