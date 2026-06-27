package value_test

import (
	"testing"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/google/go-cmp/cmp"
)

func TestChannelOfValueClose(t *testing.T) {
	opts := comparer.Options()

	ch := value.NewChannelOfValue(0)
	got := ch.Close()
	if !got.IsUndefined() {
		t.Logf("got an error when closing an open channel: %s", got.Inspect())
		t.Fail()
	}

	ch = value.NewChannelOfValue(0)
	ch.Close()
	got = ch.Close()
	want := value.Ref(value.NewError(value.ChannelClosedErrorClass, "cannot close a closed channel"))
	if diff := cmp.Diff(want, got, opts...); diff != "" {
		t.Log(diff)
		t.Fail()
	}
}

func TestChannelOfValuePush(t *testing.T) {
	opts := comparer.Options()

	ch := value.NewChannelOfValue(2)
	got := ch.Push(value.True.ToValue())
	if !got.IsUndefined() {
		t.Logf("got an error when pushing to an open channel: %s", got.Inspect())
		t.Fail()
	}

	got, gotErr := ch.Pop()
	if gotErr.IsNotUndefined() {
		t.Log("got error when popping from an open channel")
		t.Fail()
	}
	want := value.True.ToValue()
	if diff := cmp.Diff(want, got, opts...); diff != "" {
		t.Log(diff)
		t.Fail()
	}

	ch = value.NewChannelOfValue(2)
	ch.Close()
	got = ch.Push(value.Nil)
	want = value.Ref(value.NewError(value.ChannelClosedErrorClass, "cannot push values to a closed channel"))
	if diff := cmp.Diff(want, got, opts...); diff != "" {
		t.Log(diff)
		t.Fail()
	}
}

func TestChannelOfValuePop(t *testing.T) {
	opts := comparer.Options()

	ch := value.NewChannelOfValue(2)
	ch.Push(value.SmallInt(5).ToValue())

	got, gotErr := ch.Pop()
	if gotErr.IsNotUndefined() {
		t.Log("got err when popping from an open channel")
		t.Fail()
	}
	want := value.SmallInt(5).ToValue()
	if diff := cmp.Diff(want, got, opts...); diff != "" {
		t.Log(diff)
		t.Fail()
	}

	ch = value.NewChannelOfValue(2)
	ch.Close()
	_, gotErr = ch.Pop()
	if gotErr.IsUndefined() {
		t.Log("got no error when popping from a closed channel")
		t.Fail()
	}
}

func TestChannelOfValueNext(t *testing.T) {
	opts := comparer.Options()

	ch := value.NewChannelOfValue(2)
	ch.Push(value.SmallInt(5).ToValue())

	got, gotErr := ch.NextValue()
	wantErr := value.Undefined
	if diff := cmp.Diff(wantErr, gotErr, opts...); diff != "" {
		t.Log(diff)
		t.Fail()
	}
	want := value.SmallInt(5).ToValue()
	if diff := cmp.Diff(want, got, opts...); diff != "" {
		t.Log(diff)
		t.Fail()
	}

	ch = value.NewChannelOfValue(2)
	ch.Close()
	_, gotErr = ch.NextValue()
	wantErr = symbol.L_stop_iteration.ToValue()
	if diff := cmp.Diff(wantErr, gotErr, opts...); diff != "" {
		t.Log(diff)
		t.Fail()
	}
}

func TestChannelOfValueLength(t *testing.T) {
	ch := value.NewChannelOfValue(0)
	if got := ch.Length(); got != 0 {
		t.Logf("got wrong length for an unbuffered channel: %d", got)
		t.Fail()
	}

	ch = value.NewChannelOfValue(2)
	if got := ch.Length(); got != 0 {
		t.Logf("got wrong length for an empty channel: %d", got)
		t.Fail()
	}

	ch = value.NewChannelOfValue(2)
	ch.Push(value.SmallInt(5).ToValue())
	if got := ch.Length(); got != 1 {
		t.Logf("got wrong length for channel with 1 element: %d", got)
		t.Fail()
	}
}

func TestChannelOfValueCapacity(t *testing.T) {
	ch := value.NewChannelOfValue(0)
	if got := ch.Capacity(); got != 0 {
		t.Logf("got wrong cap for an unbuffered channel: %d", got)
		t.Fail()
	}

	ch = value.NewChannelOfValue(2)
	if got := ch.Capacity(); got != 2 {
		t.Logf("got wrong length for channel with 2 slots: %d", got)
		t.Fail()
	}

	ch = value.NewChannelOfValue(5)
	ch.Push(value.SmallInt(5).ToValue())
	if got := ch.Capacity(); got != 5 {
		t.Logf("got wrong length for channel with 5 slots: %d", got)
		t.Fail()
	}
}
