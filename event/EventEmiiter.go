package event_emitter

import (
	"sync"
)

// EventData es el tipo de datos que se enviará en los eventos.
type EventData struct {
	Type string
	Data interface{}
}

// EventEmitter es una estructura que manejará los eventos.
type EventEmitter struct {
	listeners map[string][]chan EventData
	lock      sync.Mutex
}

// NewEventEmitter crea una nueva instancia de EventEmitter.
func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		listeners: make(map[string][]chan EventData),
	}
}

// On registra un nuevo listener para un tipo de evento específico.
func (e *EventEmitter) On(eventType string, listener chan EventData) {
	e.lock.Lock()
	defer e.lock.Unlock()

	// Agregar el listener al slice de listeners para este evento.
	e.listeners[eventType] = append(e.listeners[eventType], listener)
}

// Emit emite un evento a todos los listeners registrados para ese tipo de evento.
func (e *EventEmitter) Emit(eventType string, data interface{}) {
	e.lock.Lock()
	defer e.lock.Unlock()

	// Notificar a todos los listeners registrados para este tipo de evento.
	for _, listener := range e.listeners[eventType] {
		listener <- EventData{Type: eventType, Data: data}
	}
}
