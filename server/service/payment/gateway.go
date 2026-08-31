package payment

import (
	"context"

	"mewsoproxy/server/model"
)

const (
	ChannelMock = "mock"
	ChannelEpay = "epay"
)

type CreateInput struct {
	Order     *model.Order
	Payment   *model.Payment
	NotifyURL string
	ReturnURL string
}

type CreateResult struct {
	Channel   string `json:"channel"`
	PayType   string `json:"pay_type"` // redirect | qrcode | completed
	URL       string `json:"url"`
	QRCode    string `json:"qr_code"`
	Completed bool   `json:"completed"`
}

type ChannelDTO struct {
	ID      int     `json:"id"`
	Payment string  `json:"payment"`
	Name    string  `json:"name"`
	Icon    *string `json:"icon,omitempty"`
}

type Gateway interface {
	Channel() string
	Create(ctx context.Context, in CreateInput) (*CreateResult, error)
	VerifyNotify(ctx context.Context, params map[string]string) error
}

func Resolve(p *model.Payment) (Gateway, error) {
	switch p.Payment {
	case ChannelMock:
		return NewMockGateway(), nil
	default:
		return NewEpayGateway(p)
	}
}
