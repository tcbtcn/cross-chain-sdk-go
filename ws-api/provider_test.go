package wsapi

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWsProviderConnector is a mock implementation of WsProviderConnector
type MockWsProviderConnector struct {
	mock.Mock
}

func (m *MockWsProviderConnector) Init() {
	m.Called()
}

func (m *MockWsProviderConnector) Send(message interface{}) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *MockWsProviderConnector) On(event WebSocketEvent, cb interface{}) {
	m.Called(event, cb)
}

func (m *MockWsProviderConnector) Off(event WebSocketEvent, cb interface{}) {
	m.Called(event, cb)
}

func (m *MockWsProviderConnector) OnOpen(cb OnOpenCb) {
	m.Called(cb)
}

func (m *MockWsProviderConnector) OnClose(cb OnCloseCb) {
	m.Called(cb)
}

func (m *MockWsProviderConnector) OnError(cb OnErrorCb) {
	m.Called(cb)
}

func (m *MockWsProviderConnector) OnMessage(cb OnMessageCb) {
	m.Called(cb)
}

func (m *MockWsProviderConnector) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockWsProviderConnector) IsConnected() bool {
	args := m.Called()
	return args.Bool(0)
}

func TestWsApiConfig(t *testing.T) {
	config := WsApiConfig{
		URL:      "wss://api.1inch.dev/fusion-plus/ws",
		AuthKey:  "test-key",
		LazyInit: false,
	}

	assert.Equal(t, "wss://api.1inch.dev/fusion-plus/ws", config.URL)
	assert.Equal(t, "test-key", config.AuthKey)
	assert.False(t, config.LazyInit)
}

func TestMockWsProviderConnector(t *testing.T) {
	// Test Init
	t.Run("Init", func(t *testing.T) {
		mockProvider := new(MockWsProviderConnector)
		mockProvider.Mock.On("Init")
		mockProvider.Init()
		mockProvider.AssertCalled(t, "Init")
	})

	// Test Send
	t.Run("Send", func(t *testing.T) {
		mockProvider := new(MockWsProviderConnector)
		mockProvider.Mock.On("Send", mock.Anything).Return(nil)
		err := mockProvider.Send(map[string]string{"test": "data"})
		assert.NoError(t, err)
		mockProvider.AssertCalled(t, "Send", mock.Anything)
	})

	// Test Send with error
	t.Run("Send_Error", func(t *testing.T) {
		mockProvider := new(MockWsProviderConnector)
		mockProvider.Mock.On("Send", mock.Anything).Return(errors.New("send error"))
		err := mockProvider.Send(map[string]string{"test": "data"})
		assert.Error(t, err)
		assert.Equal(t, "send error", err.Error())
	})

	// Test On
	t.Run("On", func(t *testing.T) {
		mockProvider := new(MockWsProviderConnector)
		cb := func() {}
		mockProvider.Mock.On("On", EventOpen, mock.AnythingOfType("func()"))
		mockProvider.On(EventOpen, cb)
		mockProvider.AssertCalled(t, "On", EventOpen, mock.AnythingOfType("func()"))
	})

	// Test Off
	t.Run("Off", func(t *testing.T) {
		mockProvider := new(MockWsProviderConnector)
		cb := func() {}
		mockProvider.Mock.On("Off", EventOpen, mock.AnythingOfType("func()"))
		mockProvider.Off(EventOpen, cb)
		mockProvider.AssertCalled(t, "Off", EventOpen, mock.AnythingOfType("func()"))
	})

	// Test OnOpen
	t.Run("OnOpen", func(t *testing.T) {
		mockProvider := new(MockWsProviderConnector)
		openCb := func() {}
		mockProvider.Mock.On("OnOpen", mock.AnythingOfType("wsapi.OnOpenCb"))
		mockProvider.OnOpen(openCb)
		mockProvider.AssertCalled(t, "OnOpen", mock.AnythingOfType("wsapi.OnOpenCb"))
	})

	// Test OnClose
	t.Run("OnClose", func(t *testing.T) {
		mockProvider := new(MockWsProviderConnector)
		closeCb := func() {}
		mockProvider.Mock.On("OnClose", mock.AnythingOfType("wsapi.OnCloseCb"))
		mockProvider.OnClose(closeCb)
		mockProvider.AssertCalled(t, "OnClose", mock.AnythingOfType("wsapi.OnCloseCb"))
	})

	// Test OnError
	t.Run("OnError", func(t *testing.T) {
		mockProvider := new(MockWsProviderConnector)
		errorCb := func(err error) {}
		mockProvider.Mock.On("OnError", mock.AnythingOfType("wsapi.OnErrorCb"))
		mockProvider.OnError(errorCb)
		mockProvider.AssertCalled(t, "OnError", mock.AnythingOfType("wsapi.OnErrorCb"))
	})

	// Test OnMessage
	t.Run("OnMessage", func(t *testing.T) {
		mockProvider := new(MockWsProviderConnector)
		msgCb := func(data interface{}) {}
		mockProvider.Mock.On("OnMessage", mock.AnythingOfType("wsapi.OnMessageCb"))
		mockProvider.OnMessage(msgCb)
		mockProvider.AssertCalled(t, "OnMessage", mock.AnythingOfType("wsapi.OnMessageCb"))
	})

	// Test Close
	t.Run("Close", func(t *testing.T) {
		mockProvider := new(MockWsProviderConnector)
		mockProvider.Mock.On("Close").Return(nil)
		err := mockProvider.Close()
		assert.NoError(t, err)
		mockProvider.AssertCalled(t, "Close")
	})

	// Test IsConnected
	t.Run("IsConnected", func(t *testing.T) {
		mockProvider := new(MockWsProviderConnector)
		mockProvider.Mock.On("IsConnected").Return(true)
		connected := mockProvider.IsConnected()
		assert.True(t, connected)
		mockProvider.AssertCalled(t, "IsConnected")
	})
}
