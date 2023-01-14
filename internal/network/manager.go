package network

import (
	"context"
	"errors"
	"github.com/mbpolan/lull/internal/state"
	"net/http"
	"sync"
)

type RequestHandler func(item *state.CollectionItem, res *http.Response, err error)

// Manager handles sending, queueing and cancelling in-flight HTTP requests.
type Manager struct {
	client      *Client
	ctx         context.Context
	cancelFunc  context.CancelFunc
	currentItem *state.CollectionItem
	handler     RequestHandler
	mutex       sync.Mutex
	pending     bool
}

// NewNetworkManager returns a new instance of Manager with the given handler function. The handler will be invoked
// whenever a network request completes, whether successful or not.
func NewNetworkManager(handler RequestHandler) *Manager {
	m := new(Manager)
	m.client = NewClient()
	m.handler = handler
	m.mutex = sync.Mutex{}
	m.pending = false

	return m
}

// Pending returns whether a request is already in-flight.
func (m *Manager) Pending() bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.pending
}

// SendRequest dispatches an HTTP request for the given collection item. An error will be returned if a request is
// already in progress.
func (m *Manager) SendRequest(item *state.CollectionItem) error {
	m.mutex.Lock()

	if m.pending {
		m.mutex.Unlock()
		return errors.New("request in progress")
	}

	m.currentItem = item
	m.pending = true
	m.ctx, m.cancelFunc = context.WithCancel(context.Background())

	// release the lock since subsequent calls to this method will error out
	m.mutex.Unlock()

	go func() {
		res, err := m.client.Exchange(m.ctx, item)
		m.handler(m.currentItem, res, err)

		m.resetCurrent()
	}()

	return nil
}

// CancelCurrent aborts the currently in-flight HTTP request.
func (m *Manager) CancelCurrent() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if !m.pending {
		return
	}

	m.cancelFunc()
	m.pending = false
}

// resetCurrent cancels the currently in-flight request and resets bookkeeping state.
func (m *Manager) resetCurrent() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.currentItem = nil
	m.pending = false
}