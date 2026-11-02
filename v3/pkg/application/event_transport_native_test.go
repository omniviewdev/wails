//go:build !server

package application

import (
	"fmt"
	"testing"
)

func TestNativeDefaultTransportPreservesPreReadyEvents(t *testing.T) {
	a := newEventTransportTestApp(t, Options{})
	probe := &threadProbeApp{}
	probe.onMain.Store(true)
	a.impl = probe
	impl := &stubWindowImpl{}
	window := &WebviewWindow{id: 1, impl: impl}
	a.Window.Add(window)

	first := &CustomEvent{Name: "before-ready-first", Data: "first payload"}
	second := &CustomEvent{Name: "before-ready-second", Data: "second payload"}
	a.Event.dispatch(first)
	a.Event.dispatch(second)

	window.pendingJSMutex.Lock()
	pending := append([]string(nil), window.pendingJS...)
	ready := window.runtimeLoaded
	window.pendingJSMutex.Unlock()
	if ready {
		t.Fatal("test window unexpectedly became runtime-ready")
	}
	if len(pending) != 2 {
		t.Fatalf("native pre-ready queue contains %d events, want two", len(pending))
	}
	for index, event := range []*CustomEvent{first, second} {
		if want := fmt.Sprintf(inlineEventJS, event.ToJSON()); pending[index] != want {
			t.Fatalf("native pre-ready event %d changed payload or order", index)
		}
	}
	if len(impl.delivered()) != 0 {
		t.Fatal("native transport executed JavaScript before runtime readiness")
	}
}
