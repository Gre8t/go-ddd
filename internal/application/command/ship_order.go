package command

import (
	"context"

	"ordercontext/internal/application/event"
	"ordercontext/internal/domain"

	"github.com/eyazici90/go-mediator/mediator"
)

type ShipOrderCommand struct {
	OrderID string `validate:"required,min=10"`
}

func (ShipOrderCommand) Key() string { return "ShipOrderCommand" }

type ShipOrderCommandHandler struct {
	repository     domain.OrderRepository
	eventPublisher event.Publisher
}

func NewShipOrderCommandHandler(r domain.OrderRepository, e event.Publisher) ShipOrderCommandHandler {
	return ShipOrderCommandHandler{
		repository:     r,
		eventPublisher: e,
	}
}

func (h ShipOrderCommandHandler) Handle(ctx context.Context, msg mediator.Message) error {
	cmd := msg.(ShipOrderCommand)
	o, err := h.repository.Get(ctx, cmd.OrderID)
	if err != nil {
		return err
	}

	err = o.Ship()

	if err != nil {
		return err
	}

	if err := h.repository.Update(ctx, o); err != nil {
		return err
	}

	h.eventPublisher.PublishAll(o.Events())

	return nil
}
