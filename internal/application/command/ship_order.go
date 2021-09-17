package command

import (
	"context"

	"ordercontext/internal/application/event"
	"ordercontext/internal/domain/order"

	"github.com/eyazici90/go-mediator/pkg/mediator"
)

type ShipOrder struct {
	OrderID string `validate:"required,min=10"`
}

func (ShipOrder) Key() string { return "ShipOrder" }

type ShipOrderHandler struct {
	orderHandler
	eventPublisher event.Publisher
}

func NewShipOrderHandler(getOrder GetOrder,
	updateOrder UpdateOrder,
	e event.Publisher) ShipOrderHandler {
	return ShipOrderHandler{
		orderHandler:   newOrderHandler(getOrder, updateOrder),
		eventPublisher: e,
	}
}

func (h ShipOrderHandler) Handle(ctx context.Context, msg mediator.Message) error {
	cmd, ok := msg.(ShipOrder)
	if err := checkType(ok); err != nil {
		return err
	}

	var ord *order.Order
	if err := h.updateErr(ctx, cmd.OrderID, func(o *order.Order) error {
		ord = o
		return ord.Ship()
	}); err != nil {
		return err
	}

	h.eventPublisher.PublishAll(ord.Events())

	return nil
}
