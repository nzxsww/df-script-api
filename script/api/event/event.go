package event

import (
	"reflect"
	"sync"
)

type Listener interface{}

type EventExecutor func(event Event)

type Priority int

const (
	PriorityLowest Priority = iota
	PriorityLow
	PriorityNormal
	PriorityHigh
	PriorityHighest
	PriorityMonitor
)

func (p Priority) Order() int {
	return int(p)
}

type Event interface {
	GetEventName() string
	GetHandlers() *HandlerList
	IsCancelled() bool
	SetCancelled(cancelled bool)
}

type HandlerList struct {
	handlers    []RegisteredListener
	handlerType reflect.Type
	lock        bool
	mu          sync.RWMutex
}

func NewHandlerList(eventType reflect.Type) *HandlerList {
	return &HandlerList{
		handlerType: eventType,
	}
}

func (h *HandlerList) Register(listener RegisteredListener) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.lock {
		panic("Cannot register new listener while handling events")
	}

	insertIndex := len(h.handlers)
	for i, existing := range h.handlers {
		if listener.Priority.Order() < existing.Priority.Order() {
			insertIndex = i
			break
		}
	}

	h.handlers = append(h.handlers, listener)
	if insertIndex < len(h.handlers)-1 {
		copy(h.handlers[insertIndex+1:], h.handlers[insertIndex:len(h.handlers)-1])
		h.handlers[insertIndex] = listener
	}
}

func (h *HandlerList) Unregister(listener RegisteredListener) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.lock {
		panic("Cannot unregister listener while handling events")
	}

	for i := range h.handlers {
		if h.handlers[i].Listener == listener.Listener {
			h.handlers = append(h.handlers[:i], h.handlers[i+1:]...)
			return
		}
	}
}

func (h *HandlerList) GetListeners() []RegisteredListener {
	h.mu.RLock()
	defer h.mu.RUnlock()
	result := make([]RegisteredListener, len(h.handlers))
	copy(result, h.handlers)
	return result
}

func (h *HandlerList) Bake() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.lock = true
}

func (h *HandlerList) HandlerType() reflect.Type {
	return h.handlerType
}

type RegisteredListener struct {
	Listener Listener
	Executor EventExecutor
	Priority Priority
	Plugin   interface{}
}

func (r RegisteredListener) GetPlugin() interface{} {
	return r.Plugin
}

func (r RegisteredListener) Execute(event Event) {
	r.Executor(event)
}

type EventHandler struct {
	Priority        Priority
	IgnoreCancelled bool
}

func (e EventHandler) GetPriority() Priority {
	if e.Priority == 0 {
		return PriorityNormal
	}
	return e.Priority
}
