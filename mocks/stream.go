// Code generated by MockGen. DO NOT EDIT.
// Source: ../internal/rep/stream.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "price_aggregator/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPriceStreamSubscriber is a mock of PriceStreamSubscriber interface.
type MockPriceStreamSubscriber struct {
	ctrl     *gomock.Controller
	recorder *MockPriceStreamSubscriberMockRecorder
}

// MockPriceStreamSubscriberMockRecorder is the mock recorder for MockPriceStreamSubscriber.
type MockPriceStreamSubscriberMockRecorder struct {
	mock *MockPriceStreamSubscriber
}

// NewMockPriceStreamSubscriber creates a new mock instance.
func NewMockPriceStreamSubscriber(ctrl *gomock.Controller) *MockPriceStreamSubscriber {
	mock := &MockPriceStreamSubscriber{ctrl: ctrl}
	mock.recorder = &MockPriceStreamSubscriberMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPriceStreamSubscriber) EXPECT() *MockPriceStreamSubscriberMockRecorder {
	return m.recorder
}

// SubscribePriceStream mocks base method.
func (m *MockPriceStreamSubscriber) SubscribePriceStream(arg0 models.Ticker) (chan models.TickerPrice, chan error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribePriceStream", arg0)
	ret0, _ := ret[0].(chan models.TickerPrice)
	ret1, _ := ret[1].(chan error)
	return ret0, ret1
}

// SubscribePriceStream indicates an expected call of SubscribePriceStream.
func (mr *MockPriceStreamSubscriberMockRecorder) SubscribePriceStream(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribePriceStream", reflect.TypeOf((*MockPriceStreamSubscriber)(nil).SubscribePriceStream), arg0)
}