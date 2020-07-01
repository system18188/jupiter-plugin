package dbr

import (
	"github.com/douyu/jupiter/pkg/xlog"
)

var JupiterReceiver = &jupiterEventReceiver{}

// jupiterEventReceiver is a sentinel EventReceiver.
type jupiterEventReceiver struct{}

// Event receives a simple notification when various events occur.
// 事件在各种事件发生时接收一个简单的通知。
func (n *jupiterEventReceiver) Event(eventName string) {
	xlog.Debug("Sql EventKv", xlog.String("EventName", eventName))
}

// EventKv receives a notification when various events occur along with
// optional key/value data.
func (n *jupiterEventReceiver) EventKv(eventName string, kvs map[string]string) {
	xlog.Debug("Sql EventKv", xlog.String("EventName", eventName), xlog.Any("kvs", kvs))
}

// EventErr receives a notification of an error if one occurs.
func (n *jupiterEventReceiver) EventErr(eventName string, err error) error {
	xlog.Error("Sql EventErr", xlog.FieldMod("dbr"), xlog.FieldName(eventName), xlog.FieldErr(err))
	return err
}

// EventErrKv receives a notification of an error if one occurs along with
// optional key/value data.
func (n *jupiterEventReceiver) EventErrKv(eventName string, err error, kvs map[string]string) error {
	xlog.Error("Sql EventErrKv", xlog.FieldMod("dbr"), xlog.FieldName(eventName), xlog.FieldErr(err), xlog.FieldValueAny(kvs))
	return err
}

// Timing receives the time an event took to happen.
func (n *jupiterEventReceiver) Timing(eventName string, nanoseconds int64) {
	xlog.Debug("Sql Timing", xlog.FieldMod("dbr"), xlog.FieldName(eventName), xlog.Any("exec time", nanoseconds))
}

// TimingKv receives the time an event took to happen along with optional key/value data.
func (n *jupiterEventReceiver) TimingKv(eventName string, nanoseconds int64, kvs map[string]string) {
	xlog.Debug("Sql TimingKv", xlog.FieldMod("dbr"), xlog.FieldName(eventName), xlog.Any("exec time", nanoseconds), xlog.FieldValueAny(kvs))
}
