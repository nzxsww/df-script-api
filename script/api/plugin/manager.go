package plugin

import (
	"reflect"

	"github.com/df-mc/dragonfly/server"
	"github.com/nzxsww/dragonfly-script-api/script/api/event"
)

var globalManager *Manager

type Manager struct {
	srv     *server.Server
	plugins map[string]Plugin
}

func NewManager(srv *server.Server) *Manager {
	globalManager = &Manager{
		srv:     srv,
		plugins: make(map[string]Plugin),
	}
	return globalManager
}

func GetManager() *Manager {
	return globalManager
}

func (m *Manager) Server() *server.Server {
	return m.srv
}

func (m *Manager) RegisterEvents(listener event.Listener, p Plugin) {
	t := reflect.TypeOf(listener)
	v := reflect.ValueOf(listener)

	if t.Kind() != reflect.Ptr {
		t = reflect.PtrTo(t)
		v = v.Addr()
	}

	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		if method.Name == "GetHandlerList" || method.Name == "GetHandlers" {
			continue
		}

		if method.Type.NumIn() != 1 {
			continue
		}

		eventType := method.Type.In(0)
		if !eventType.Implements(eventInterface) {
			continue
		}

		eventName := getEventName(eventType)
		handlers := event.GetEventHandlerList(eventName)
		if handlers == nil {
			continue
		}

		executor := func(e event.Event) {
			methodVal := v.MethodByName(method.Name)
			if methodVal.IsValid() {
				eventVal := reflect.ValueOf(e)
				if eventVal.Type() == eventType {
					methodVal.Call([]reflect.Value{eventVal})
				} else if eventVal.Type().Implements(eventType) {
					converted := reflect.Zero(eventType)
					if converted.CanSet() {
						converted.Set(eventVal)
						methodVal.Call([]reflect.Value{converted})
					}
				}
			}
		}

		listener := event.RegisteredListener{
			Listener: listener,
			Executor: executor,
			Priority: event.PriorityNormal,
			Plugin:   p,
		}

		handlers.Register(listener)
	}
}

func (m *Manager) CallEvent(e event.Event) {
	event.CallEvent(e)
}

func (m *Manager) GetPlugin(name string) Plugin {
	return m.plugins[name]
}

func (m *Manager) GetPlugins() []Plugin {
	result := make([]Plugin, 0, len(m.plugins))
	for _, p := range m.plugins {
		result = append(result, p)
	}
	return result
}

func (m *Manager) AddPlugin(p Plugin) {
	m.plugins[p.GetName()] = p
}

var eventInterface = reflect.TypeOf((*event.Event)(nil)).Elem()

func getEventName(t reflect.Type) string {
	name := t.Name()
	if name == "" {
		return t.String()
	}
	return name
}
