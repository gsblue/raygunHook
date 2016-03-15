package raygunHook

import (
	"errors"
	"net/http"
	"testing"

	"reflect"

	"github.com/Sirupsen/logrus"
	"github.com/gsblue/raygunclient"
	"github.com/gsblue/raygunclient/stack"
)

func TestNewHook_WhenValidConfig_ShouldReturnHook(t *testing.T) {
	h := NewHook(&HookConfig{APIKey: "key"})

	if h == nil {
		t.Error("expected hook to be not nil")
	}
}

func TestFire_WhenLogEntryHasError_ShouldSubmitEntry(t *testing.T) {
	e := errors.New("test error")
	r, _ := http.NewRequest("GET", "/", nil)
	u := "user"
	cd := &struct{ name string }{"custom"}

	h := &hook{
		notifier: &raygunClientMock{
			mockNotify: func(entry *raygunclient.ErrorEntry) error {
				if entry == nil {
					t.Error("entry should not be nil")
				}
				return nil
			},
		},
	}

	le := logrus.NewEntry(logrus.StandardLogger()).
		WithField(ErrorFieldName, e)

	le = EntryWithRequest(le, r)
	le = EntryWithUser(le, u)
	le = EntryWithCustomData(le, cd)

	if err := h.Fire(le); err != nil {
		t.Error(err)
	}
}

func TestFire_WhenLogEntryHasErrorWithStackTrace_ShouldSubmitEntryWithStackTrace(t *testing.T) {
	e := errors.New("test error")
	trace := make(stack.Trace, 0, 0)
	trace.AddEntry(23, "github.com/gsblue/raygunclient", "main.go", "Notify")

	h := &hook{
		notifier: &raygunClientMock{
			mockNotifyWithStackTrace: func(entry *raygunclient.ErrorEntry, s stack.Trace) error {
				if entry == nil {
					t.Error("entry should not be nil")
				} else if s == nil {
					t.Error("trace should not be nil")
				} else if !reflect.DeepEqual(s, trace) {
					t.Errorf("expected %v, got %v", trace, s)
				}
				return nil
			},
		},
	}

	le := logrus.NewEntry(logrus.StandardLogger()).
		WithField(ErrorFieldName, e)

	le = EntryWithStackTrace(le, trace)

	if err := h.Fire(le); err != nil {
		t.Error(err)
	}
}

func TestFire_WhenLogEntryErrorMsg_ShouldSubmitEntry(t *testing.T) {
	msg := "test error"
	r, _ := http.NewRequest("GET", "/", nil)
	u := "user"
	cd := &struct{ name string }{"custom"}

	h := &hook{
		notifier: &raygunClientMock{
			mockNotify: func(entry *raygunclient.ErrorEntry) error {
				if entry == nil {
					t.Error("entry should not be nil")
				}
				return nil
			},
		},
	}

	le := logrus.NewEntry(logrus.StandardLogger())
	le = EntryWithRequest(le, r)
	le = EntryWithUser(le, u)
	le = EntryWithCustomData(le, cd)
	le.Message = msg

	if err := h.Fire(le); err != nil {
		t.Error(err)
	}
}

func TestFire_WhenClientReturnsError_ShouldReturnError(t *testing.T) {
	errClient := errors.New("client error")
	h := &hook{
		notifier: &raygunClientMock{
			mockNotify: func(entry *raygunclient.ErrorEntry) error {
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
	mockNotify               notifyFn
	mockNotifyWithStackTrace notifyWithStackTraceFn
}

func (r *raygunClientMock) Notify(entry *raygunclient.ErrorEntry) error {
	return r.mockNotify(entry)
}

func (r *raygunClientMock) NotifyWithStackTrace(entry *raygunclient.ErrorEntry, s stack.Trace) error {
	return r.mockNotifyWithStackTrace(entry, s)
}

type notifyFn func(entry *raygunclient.ErrorEntry) error
type notifyWithStackTraceFn func(entry *raygunclient.ErrorEntry, s stack.Trace) error
