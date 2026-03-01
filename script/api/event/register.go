package event

import "sync"

var (
	registeredEvents = make(map[string]*HandlerList)
	eventsMu         sync.RWMutex
)

func RegisterEvent(name string, handlers *HandlerList) {
	eventsMu.Lock()
	defer eventsMu.Unlock()
	registeredEvents[name] = handlers
}

func GetEventHandlerList(eventName string) *HandlerList {
	eventsMu.RLock()
	defer eventsMu.RUnlock()
	return registeredEvents[eventName]
}

func CallEvent(event Event) {
	name := event.GetEventName()
	eventsMu.RLock()
	handlers := registeredEvents[name]
	eventsMu.RUnlock()

	if handlers == nil {
		return
	}

	listeners := handlers.GetListeners()
	sorted := make([]RegisteredListener, len(listeners))
	copy(sorted, listeners)

	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].Priority.Order() < sorted[i].Priority.Order() {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	for _, l := range sorted {
		l.Executor(event)
	}
}

func BakeAllEvents() {
	eventsMu.RLock()
	defer eventsMu.RUnlock()
	for _, hl := range registeredEvents {
		hl.Bake()
	}
}
