package event

import "time"

type BalanceUpdated struct {
	Name    string
	Payload any
}

func NewBalanceUpdatedEvent() *BalanceUpdated {
	return &BalanceUpdated{
		Name: "BalanceUpdated",
	}
}

func (e *BalanceUpdated) GetName() string {
	return e.Name
}

func (e *BalanceUpdated) GetPayload() any {
	return e.Payload
}

func (e *BalanceUpdated) SetPayload(payload any) {
	e.Payload = payload
}

func (e *BalanceUpdated) GetDateTime() time.Time {
	return time.Now()
}
