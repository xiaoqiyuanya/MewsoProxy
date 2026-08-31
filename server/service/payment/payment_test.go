package payment

import (
	"context"
	"testing"

	"mewsoproxy/server/model"
)

func intp(v int) *int { return &v }

func TestEpaySign(t *testing.T) {
	params := map[string]string{
		"pid":          "1001",
		"type":         "alipay",
		"out_trade_no": "ORD-abc",
		"notify_url":   "http://x/api/v1/payment/notify",
		"return_url":   "http://x/#/order",
		"name":         "套餐订单",
		"money":        "9.90",
	}
	got := epaySign(params, "demo-key")
	if got == "" {
		t.Fatal("sign empty")
	}
	want := "7634ebd9b5802079360ff1e704554809"
	if got != want {
		t.Fatalf("sign = %s, want %s", got, want)
	}
}

func TestMockGatewayCreateAndVerify(t *testing.T) {
	gw := NewMockGateway()
	order := &model.Order{TradeNo: "ORD-mock-1", TotalAmount: 9900}
	res, err := gw.Create(context.Background(), CreateInput{Order: order, NotifyURL: "http://x/api/v1/payment/notify"})
	if err != nil {
		t.Fatalf("create error: %v", err)
	}
	if res.Channel != ChannelMock || res.PayType != "qrcode" {
		t.Fatalf("unexpected result: %+v", res)
	}
	if res.QRCode == "" {
		t.Fatal("qr_code empty")
	}
	if err := gw.VerifyNotify(context.Background(), map[string]string{"trade_no": order.TradeNo}); err != nil {
		t.Fatalf("verify error: %v", err)
	}
	if err := gw.VerifyNotify(context.Background(), map[string]string{}); err == nil {
		t.Fatal("verify should fail with missing trade_no")
	}
}

func TestResolveMock(t *testing.T) {
	p := &model.Payment{Payment: ChannelMock, Config: ""}
	gw, err := Resolve(p)
	if err != nil {
		t.Fatalf("resolve error: %v", err)
	}
	if gw.Channel() != ChannelMock {
		t.Fatalf("channel = %s, want %s", gw.Channel(), ChannelMock)
	}
}
