package bbgo

import (
	"sync"

	"github.com/c9s/bbgo/pkg/types"
)

type OrderStore struct {
	// any created orders for tracking trades
	mu     sync.Mutex
	orders map[uint64]types.Order

	RemoveCancelled bool
}

func NewOrderStore() *OrderStore {
	return &OrderStore{
		orders: make(map[uint64]types.Order),
	}
}

func (s *OrderStore) Exists(oID uint64) (ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok = s.orders[oID]
	return ok
}

func (s *OrderStore) Add(orders ...types.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, o := range orders {
		s.orders[o.OrderID] = o
	}
}

func (s *OrderStore) Remove(o types.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.orders, o.OrderID)
}

func (s *OrderStore) Update(o types.Order) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.orders[o.OrderID]
	if ok {
		s.orders[o.OrderID] = o
	}
	return ok
}

func (s *OrderStore) BindStream(stream types.Stream) {
	stream.OnOrderUpdate(s.handleOrderUpdate)
}

func (s *OrderStore) handleOrderUpdate(order types.Order) {
	switch order.Status {
	case types.OrderStatusPartiallyFilled, types.OrderStatusNew, types.OrderStatusFilled:
		s.Update(order)

	case types.OrderStatusCanceled:
		if s.RemoveCancelled {
			s.Remove(order)
		}

	case types.OrderStatusRejected:
		s.Remove(order)
	}
}
