package event_test

import (
	"reflect"
	"sync"
	"testing"

	"github.com/nzxsww/dragonfly-script-api/script/api/event"
)

// --- Evento de prueba minimal (no necesita Dragonfly) ---

type mockHandlerList struct {
	hl *event.HandlerList
}

type testEvent struct {
	name      string
	cancelled bool
	handlers  *event.HandlerList
}

func newTestEvent(name string) *testEvent {
	hl := event.NewHandlerList(reflect.TypeOf(&testEvent{}))
	e := &testEvent{name: name, handlers: hl}
	return e
}

func (e *testEvent) GetEventName() string        { return e.name }
func (e *testEvent) GetHandlers() *event.HandlerList { return e.handlers }
func (e *testEvent) IsCancelled() bool           { return e.cancelled }
func (e *testEvent) SetCancelled(c bool)         { e.cancelled = c }

// --- Tests de HandlerList ---

func TestHandlerList_RegisterAndExecute(t *testing.T) {
	hl := event.NewHandlerList(reflect.TypeOf(&testEvent{}))
	called := false

	hl.Register(event.RegisteredListener{
		Listener: struct{}{},
		Executor: func(e event.Event) { called = true },
		Priority: event.PriorityNormal,
	})

	ev := newTestEvent("test")
	for _, l := range hl.GetListeners() {
		l.Execute(ev)
	}

	if !called {
		t.Error("el executor no fue llamado")
	}
}

func TestHandlerList_PriorityOrder(t *testing.T) {
	hl := event.NewHandlerList(reflect.TypeOf(&testEvent{}))
	order := []string{}

	// Registrar en orden mezclado, deben ejecutarse por prioridad
	hl.Register(event.RegisteredListener{
		Listener: struct{}{},
		Executor: func(e event.Event) { order = append(order, "HIGH") },
		Priority: event.PriorityHigh,
	})
	hl.Register(event.RegisteredListener{
		Listener: struct{}{},
		Executor: func(e event.Event) { order = append(order, "LOWEST") },
		Priority: event.PriorityLowest,
	})
	hl.Register(event.RegisteredListener{
		Listener: struct{}{},
		Executor: func(e event.Event) { order = append(order, "NORMAL") },
		Priority: event.PriorityNormal,
	})

	// CallEvent ordena por prioridad internamente
	ev := newTestEvent("priority_test")
	event.RegisterEvent("priority_test", hl)
	event.CallEvent(ev)

	if len(order) != 3 {
		t.Fatalf("se esperaban 3 llamadas, got %d", len(order))
	}
	if order[0] != "LOWEST" {
		t.Errorf("primero debería ser LOWEST, got %s", order[0])
	}
	if order[1] != "NORMAL" {
		t.Errorf("segundo debería ser NORMAL, got %s", order[1])
	}
	if order[2] != "HIGH" {
		t.Errorf("tercero debería ser HIGH, got %s", order[2])
	}
}

func TestHandlerList_Unregister(t *testing.T) {
	hl := event.NewHandlerList(reflect.TypeOf(&testEvent{}))
	called := false

	listener := event.RegisteredListener{
		Listener: "marker",
		Executor: func(e event.Event) { called = true },
		Priority: event.PriorityNormal,
	}
	hl.Register(listener)
	hl.Unregister(listener)

	for _, l := range hl.GetListeners() {
		l.Execute(newTestEvent("test"))
	}

	if called {
		t.Error("el executor fue llamado después de unregister")
	}
}

func TestHandlerList_BakePreventsRegister(t *testing.T) {
	hl := event.NewHandlerList(reflect.TypeOf(&testEvent{}))
	hl.Bake()

	defer func() {
		if r := recover(); r == nil {
			t.Error("se esperaba panic al registrar después de Bake()")
		}
	}()

	hl.Register(event.RegisteredListener{
		Listener: struct{}{},
		Executor: func(e event.Event) {},
		Priority: event.PriorityNormal,
	})
}

func TestHandlerList_MultipleListeners(t *testing.T) {
	hl := event.NewHandlerList(reflect.TypeOf(&testEvent{}))
	count := 0

	for i := 0; i < 5; i++ {
		hl.Register(event.RegisteredListener{
			Listener: i, // identificador único
			Executor: func(e event.Event) { count++ },
			Priority: event.PriorityNormal,
		})
	}

	for _, l := range hl.GetListeners() {
		l.Execute(newTestEvent("test"))
	}

	if count != 5 {
		t.Errorf("se esperaban 5 llamadas, got %d", count)
	}
}

func TestHandlerList_ThreadSafe(t *testing.T) {
	hl := event.NewHandlerList(reflect.TypeOf(&testEvent{}))
	var wg sync.WaitGroup

	// Registrar desde múltiples goroutines simultáneamente
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			hl.Register(event.RegisteredListener{
				Listener: id,
				Executor: func(e event.Event) {},
				Priority: event.PriorityNormal,
			})
		}(i)
	}
	wg.Wait()

	listeners := hl.GetListeners()
	if len(listeners) != 50 {
		t.Errorf("se esperaban 50 listeners, got %d", len(listeners))
	}
}

// --- Tests de cancelación ---

func TestEvent_CancelFlow(t *testing.T) {
	hl := event.NewHandlerList(reflect.TypeOf(&testEvent{}))
	event.RegisterEvent("cancel_test", hl)

	secondCalled := false

	// Primer listener cancela el evento
	hl.Register(event.RegisteredListener{
		Listener: "first",
		Executor: func(e event.Event) { e.SetCancelled(true) },
		Priority: event.PriorityLowest,
	})
	// Segundo listener solo observa (en un sistema real ignoraría si está cancelado,
	// pero aquí verificamos que el estado de cancelación persiste)
	hl.Register(event.RegisteredListener{
		Listener: "second",
		Executor: func(e event.Event) { secondCalled = e.IsCancelled() },
		Priority: event.PriorityNormal,
	})

	ev := newTestEvent("cancel_test")
	event.CallEvent(ev)

	if !ev.IsCancelled() {
		t.Error("el evento debería estar cancelado")
	}
	if !secondCalled {
		t.Error("el segundo listener debería haber visto el evento cancelado")
	}
}

// --- Tests de Priority.Order() ---

func TestPriority_Order(t *testing.T) {
	cases := []struct {
		p     event.Priority
		order int
	}{
		{event.PriorityLowest, 0},
		{event.PriorityLow, 1},
		{event.PriorityNormal, 2},
		{event.PriorityHigh, 3},
		{event.PriorityHighest, 4},
		{event.PriorityMonitor, 5},
	}
	for _, c := range cases {
		if c.p.Order() != c.order {
			t.Errorf("Priority %v: expected order %d, got %d", c.p, c.order, c.p.Order())
		}
	}
}

// --- Tests de RegisterEvent / GetEventHandlerList ---

func TestRegisterEvent_GetHandlerList(t *testing.T) {
	hl := event.NewHandlerList(reflect.TypeOf(&testEvent{}))
	event.RegisterEvent("mi_evento_unico", hl)

	got := event.GetEventHandlerList("mi_evento_unico")
	if got == nil {
		t.Fatal("GetEventHandlerList devolvió nil para evento registrado")
	}
	if got != hl {
		t.Error("GetEventHandlerList devolvió una HandlerList diferente")
	}
}

func TestGetEventHandlerList_NotFound(t *testing.T) {
	got := event.GetEventHandlerList("evento_que_no_existe_xyz")
	if got != nil {
		t.Error("GetEventHandlerList debería devolver nil para eventos no registrados")
	}
}
