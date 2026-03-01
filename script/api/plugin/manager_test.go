package plugin_test

import (
	"reflect"
	"testing"

	"github.com/nzxsww/dragonfly-script-api/script/api/event"
	"github.com/nzxsww/dragonfly-script-api/script/api/plugin"
)

// --- Evento de prueba minimal ---

type mockEvent struct {
	name      string
	cancelled bool
	handlers  *event.HandlerList
}

func newMockEvent(name string) *mockEvent {
	hl := event.NewHandlerList(reflect.TypeOf(&mockEvent{}))
	event.RegisterEvent(name, hl)
	return &mockEvent{name: name, handlers: hl}
}

func (e *mockEvent) GetEventName() string           { return e.name }
func (e *mockEvent) GetHandlers() *event.HandlerList { return e.handlers }
func (e *mockEvent) IsCancelled() bool              { return e.cancelled }
func (e *mockEvent) SetCancelled(c bool)            { e.cancelled = c }

// --- Plugin de prueba minimal ---

type mockPlugin struct {
	plugin.BasePlugin
	joinCalled bool
	chatCalled bool
}

func newMockPlugin(name string) *mockPlugin {
	return &mockPlugin{
		BasePlugin: *plugin.NewBasePlugin(name, "/tmp/"+name, nil),
	}
}

// OnMockJoin es detectado por reflexión como handler de mockJoinEvent
type mockJoinEvent struct {
	mockEvent
}

func newMockJoinEvent() *mockJoinEvent {
	hl := event.NewHandlerList(reflect.TypeOf(&mockJoinEvent{}))
	event.RegisterEvent("mockJoinEvent", hl)
	return &mockJoinEvent{mockEvent: mockEvent{name: "mockJoinEvent", handlers: hl}}
}

func (e *mockJoinEvent) GetHandlers() *event.HandlerList { return e.handlers }

// --- Tests del Manager ---

func TestManager_AddAndGetPlugin(t *testing.T) {
	mgr := plugin.NewManager(nil)

	p := newMockPlugin("TestPlugin")
	mgr.AddPlugin(p)

	got := mgr.GetPlugin("TestPlugin")
	if got == nil {
		t.Fatal("GetPlugin devolvió nil para plugin registrado")
	}
	if got.GetName() != "TestPlugin" {
		t.Errorf("expected 'TestPlugin', got '%s'", got.GetName())
	}
}

func TestManager_GetPlugins(t *testing.T) {
	mgr := plugin.NewManager(nil)

	mgr.AddPlugin(newMockPlugin("Plugin1"))
	mgr.AddPlugin(newMockPlugin("Plugin2"))
	mgr.AddPlugin(newMockPlugin("Plugin3"))

	plugins := mgr.GetPlugins()
	if len(plugins) != 3 {
		t.Errorf("expected 3 plugins, got %d", len(plugins))
	}
}

func TestManager_GetPlugin_NotFound(t *testing.T) {
	mgr := plugin.NewManager(nil)
	got := mgr.GetPlugin("NoExiste")
	if got != nil {
		t.Error("GetPlugin debería retornar nil para plugins no registrados")
	}
}

func TestManager_CallEvent_DispatchesToHandlers(t *testing.T) {
	mgr := plugin.NewManager(nil)
	called := false

	ev := newMockEvent("dispatch_test_event")
	ev.GetHandlers().Register(event.RegisteredListener{
		Listener: "test",
		Executor: func(e event.Event) { called = true },
		Priority: event.PriorityNormal,
	})

	mgr.CallEvent(ev)

	if !called {
		t.Error("CallEvent no despachó el evento al handler registrado")
	}
}

func TestManager_CallEvent_PassesCorrectEvent(t *testing.T) {
	mgr := plugin.NewManager(nil)
	var received event.Event

	ev := newMockEvent("pass_event_test")
	ev.GetHandlers().Register(event.RegisteredListener{
		Listener: "test",
		Executor: func(e event.Event) { received = e },
		Priority: event.PriorityNormal,
	})

	mgr.CallEvent(ev)

	if received == nil {
		t.Fatal("el handler no recibió ningún evento")
	}
	if received.GetEventName() != "pass_event_test" {
		t.Errorf("evento incorrecto: expected 'pass_event_test', got '%s'", received.GetEventName())
	}
}

func TestManager_CallEvent_MultipleHandlers(t *testing.T) {
	mgr := plugin.NewManager(nil)
	count := 0

	ev := newMockEvent("multi_handler_test")
	for i := 0; i < 3; i++ {
		ev.GetHandlers().Register(event.RegisteredListener{
			Listener: i,
			Executor: func(e event.Event) { count++ },
			Priority: event.PriorityNormal,
		})
	}

	mgr.CallEvent(ev)

	if count != 3 {
		t.Errorf("expected 3 handler calls, got %d", count)
	}
}

func TestManager_GetManager_Singleton(t *testing.T) {
	mgr := plugin.NewManager(nil)
	got := plugin.GetManager()
	if got != mgr {
		t.Error("GetManager() debe retornar la misma instancia que NewManager()")
	}
}

func TestBasePlugin_GetName(t *testing.T) {
	p := plugin.NewBasePlugin("MiPlugin", "/data/miplugin", nil)
	if p.GetName() != "MiPlugin" {
		t.Errorf("expected 'MiPlugin', got '%s'", p.GetName())
	}
}

func TestBasePlugin_GetDataFolder(t *testing.T) {
	p := plugin.NewBasePlugin("MiPlugin", "/data/miplugin", nil)
	if p.GetDataFolder() != "/data/miplugin" {
		t.Errorf("expected '/data/miplugin', got '%s'", p.GetDataFolder())
	}
}

func TestBasePlugin_OnEnableDisable_NoOp(t *testing.T) {
	// BasePlugin.OnEnable y OnDisable son no-ops, no deben panic
	p := plugin.NewBasePlugin("test", "/tmp", nil)
	p.OnEnable()
	p.OnDisable()
}
