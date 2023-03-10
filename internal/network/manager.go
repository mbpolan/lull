package network

import (
	"context"
	"errors"
	"github.com/mbpolan/lull/internal/state"
	"io"
	"net/http"
	"sync"
	"time"
)

type RequestHandler func(item *state.CollectionItem, result *Result)

// Manager handles sending, queueing and cancelling in-flight HTTP requests.
type Manager struct {
	client      *Client
	ctx         context.Context
	cancelFunc  context.CancelFunc
	currentItem *state.CollectionItem
	handler     RequestHandler
	mutex       sync.Mutex
	pending     bool
	startTime   time.Time
}

// Result contains the outcome of an HTTP request.
type Result struct {
	Response     *http.Response
	Payload      []byte
	PayloadError error
	Error        error
	StartTime    time.Time
	EndTime      time.Time
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
	m.startTime = time.Now()
	m.ctx, m.cancelFunc = context.WithCancel(context.Background())

	// release the lock since subsequent calls to this method will error out
	m.mutex.Unlock()

	go func() {
		authFunc, err := m.authenticate(item)
		if err != nil {
			// TODO
			panic(err)
		}

		res, err := m.client.Exchange(m.ctx, item, authFunc)

		// read the entire body and capture any errors in the process
		var payload []byte
		var payloadErr error
		if err == nil {
			payload, payloadErr = io.ReadAll(res.Body)
			defer res.Body.Close()
		}

		m.handler(m.currentItem, &Result{
			Response:     res,
			Error:        err,
			Payload:      payload,
			PayloadError: payloadErr,
			StartTime:    m.startTime,
			EndTime:      time.Now(),
		})

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

func (m *Manager) authenticate(item *state.CollectionItem) (func(req *http.Request) error, error) {
	if item.Authentication.None() {
		return nil, nil
	}

	authReq, err := item.Authentication.Data.Prepare()
	if err != nil {
		return nil, err
	}

	var res *http.Response
	if authReq != nil {
		res, err = m.client.ExchangeRequest(authReq)
		if err != nil {
			return nil, err
		}
	}

	return func(req *http.Request) error {
		return item.Authentication.Data.Apply(req, res)
	}, nil
}

// resetCurrent cancels the currently in-flight request and resets bookkeeping state.
func (m *Manager) resetCurrent() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.currentItem = nil
	m.pending = false
}
