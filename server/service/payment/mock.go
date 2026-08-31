package payment

import (
	"context"
	"strconv"

	"mewsoproxy/server/pkg/apperror"
)

type MockGateway struct{}

func NewMockGateway() *MockGateway { return &MockGateway{} }

func (g *MockGateway) Channel() string { return ChannelMock }

func (g *MockGateway) Create(ctx context.Context, in CreateInput) (*CreateResult, error) {
	return &CreateResult{
		Channel: ChannelMock,
		PayType: "qrcode",
		QRCode:  "mock://pay?trade_no=" + in.Order.TradeNo + "&amount=" + strconv.Itoa(in.Order.TotalAmount),
		URL:     in.NotifyURL,
	}, nil
}

func (g *MockGateway) VerifyNotify(ctx context.Context, params map[string]string) error {
	if params["out_trade_no"] == "" && params["trade_no"] == "" {
		return apperror.New(apperror.CodeParamFormat, "订单号缺失")
	}
	return nil
}
