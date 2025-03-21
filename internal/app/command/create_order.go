package command

import (
	"context"
	"time"

	"github.com/eyazici90/go-ddd/internal/domain"
	"github.com/eyazici90/go-ddd/pkg/aggregate"

	"github.com/eyazici90/go-mediator/mediator"
	"github.com/pkg/errors"
)

type OrderCreator interface {
	Create(context.Context, *domain.Order) error
}

type CreateOrder struct {
	ID string `validate:"required,min=10"`
}

func (CreateOrder) Key() int { return createCommandKey }

type CreateOrderHandler struct {
	orderCreator OrderCreator
}

func NewCreateOrderHandler(orderCreator OrderCreator) CreateOrderHandler {
	return CreateOrderHandler{orderCreator: orderCreator}
}

func (h CreateOrderHandler) Handle(ctx context.Context, msg mediator.Message) error {
	cmd, ok := msg.(CreateOrder)
	if err := checkType(ok); err != nil {
		return err
	}

	order, err := domain.NewOrder(domain.OrderID(cmd.ID), domain.NewCustomerID(), domain.NewProductID(), time.Now,
		domain.Submitted, aggregate.NewVersion())
	if err != nil {
		return errors.Wrap(err, "new order")
	}

	return h.orderCreator.Create(ctx, order)
}
