package raygunHook

import (
	"errors"
	"github.com/gsblue/raygunHook/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/gsblue/raygunHook/Godeps/_workspace/src/github.com/gsblue/raygun4go"
	"net/http"
	"testing"
)

func TestNewHook_WhenInvalidConfig_ShouldReturnError(t *testing.T) {
	if h, e := NewHook(&HookConfig{}); e == nil {
		t.Error("expected error from client")
	} else if h != nil {
		t.Error("expected hook to be nil")
	}
}

func TestNewHook_WhenValidConfig_ShouldReturnHook(t *testing.T) {
	if h, e := NewHook(&HookConfig{ApiKey: "key", AppName: "app name"}); e != nil {
		t.Errorf("did not expect error: %s", e)
	} else if h == nil {
		t.Error("expected hook to be not nil")
	}
}

func TestFire_WhenLogEntryHasError_ShouldSubmitEntry(t *testing.T) {
	e := errors.New("test error")
	r := &http.Request{}
	u := "user"
	cd := &struct{ name string }{"custom"}
	client, _ := raygun4go.New("key", "app name", "1.0", nil)
	h := &hook{
		Client: &raygunClientMock{
			c: client,
			fn: func(entry *raygun4go.ErrorEntry) error {
				if entry == nil {
					t.Error("entry should not be nil")
				}
				return nil
			},
		},
	}

	le := logrus.NewEntry(logrus.StandardLogger()).
		WithField(ErrorFieldName, e).
		WithField(RequestFieldName, r).
		WithField(UserFieldName, u).
		WithField(CustomDataFieldName, cd)

	if err := h.Fire(le); err != nil {
		t.Error(err)
	}
}

func TestFire_WhenLogEntryErrorMsg_ShouldSubmitEntry(t *testing.T) {
	msg := "test error"
	r := &http.Request{}
	u := "user"
	cd := &struct{ name string }{"custom"}
	client, _ := raygun4go.New("key", "app name", "1.0", nil)
	h := &hook{
		Client: &raygunClientMock{
			c: client,
			fn: func(entry *raygun4go.ErrorEntry) error {
				if entry == nil {
					t.Error("entry should not be nil")
				}
				return nil
			},
		},
	}

	le := logrus.NewEntry(logrus.StandardLogger()).
		WithField(RequestFieldName, r).
		WithField(UserFieldName, u).
		WithField(CustomDataFieldName, cd)
	le.Message = msg

	if err := h.Fire(le); err != nil {
		t.Error(err)
	}
}

func TestFire_WhenClientReturnsError_ShouldReturnError(t *testing.T) {
	errClient := errors.New("client error")
	client, _ := raygun4go.New("key", "app name", "1.0", nil)
	h := &hook{
		Client: &raygunClientMock{
			c: client,
			fn: func(entry *raygun4go.ErrorEntry) error {
				return errClient
			},
		},
	}

	le := logrus.NewEntry(logrus.StandardLogger())
	le.Message = "something went wrong"

	if err := h.Fire(le); err != errClient {
		t.Error(err)
	}
}

type raygunClientMock struct {
	c  *raygun4go.Client
	fn submitErrMockFn
}

func (r *raygunClientMock) CreateErrorEntry(err error) *raygun4go.ErrorEntry {
	return r.c.CreateErrorEntry(err)
}

func (r *raygunClientMock) CreateErrorEntryFromMsg(msg string) *raygun4go.ErrorEntry {
	return r.c.CreateErrorEntryFromMsg(msg)
}

func (r *raygunClientMock) SubmitError(entry *raygun4go.ErrorEntry) error {
	return r.fn(entry)
}

type submitErrMockFn func(entry *raygun4go.ErrorEntry) error
