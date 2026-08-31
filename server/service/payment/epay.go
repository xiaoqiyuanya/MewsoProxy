package payment

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"mewsoproxy/server/model"
	"mewsoproxy/server/pkg/apperror"
)

type epayConfig struct {
	GatewayURL string `json:"gateway_url"`
	PID        string `json:"pid"`
	Key        string `json:"key"`
	Type       string `json:"type"`
}

type EpayGateway struct {
	cfg epayConfig
}

func NewEpayGateway(p *model.Payment) (*EpayGateway, error) {
	cfg, err := parseEpayConfig(p.Config)
	if err != nil {
		return nil, err
	}
	return &EpayGateway{cfg: cfg}, nil
}

func (g *EpayGateway) Channel() string { return ChannelEpay }

func (g *EpayGateway) Create(ctx context.Context, in CreateInput) (*CreateResult, error) {
	if g.cfg.GatewayURL == "" || g.cfg.PID == "" || g.cfg.Key == "" {
		return nil, apperror.New(apperror.CodeParamFormat, "支付渠道配置不完整")
	}
	money := fmt.Sprintf("%.2f", float64(in.Order.TotalAmount)/100)
	params := map[string]string{
		"pid":          g.cfg.PID,
		"type":         g.cfg.Type,
		"out_trade_no": in.Order.TradeNo,
		"notify_url":   in.NotifyURL,
		"return_url":   in.ReturnURL,
		"name":         "套餐订单",
		"money":        money,
	}
	base := strings.TrimRight(g.cfg.GatewayURL, "/") + "/submit.php"
	u, err := url.Parse(base)
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeUpstreamError, "网关地址不合法", err)
	}
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	q.Set("sign", epaySign(params, g.cfg.Key))
	q.Set("sign_type", "MD5")
	u.RawQuery = q.Encode()
	return &CreateResult{
		Channel: ChannelEpay,
		PayType: "redirect",
		URL:     u.String(),
	}, nil
}

func (g *EpayGateway) VerifyNotify(ctx context.Context, params map[string]string) error {
	sign := params["sign"]
	copied := make(map[string]string, len(params))
	for k, v := range params {
		copied[k] = v
	}
	delete(copied, "sign")
	delete(copied, "sign_type")
	if copied["trade_status"] != "" && copied["trade_status"] != "TRADE_SUCCESS" {
		return apperror.New(apperror.CodeOrderInvalid, "支付未成功")
	}
	expected := epaySign(copied, g.cfg.Key)
	if !strings.EqualFold(expected, sign) {
		return apperror.New(apperror.CodeParamFormat, "签名校验失败")
	}
	return nil
}

func parseEpayConfig(config string) (epayConfig, error) {
	var cfg epayConfig
	if config == "" {
		return cfg, nil
	}
	if err := json.Unmarshal([]byte(config), &cfg); err != nil {
		return cfg, apperror.Wrap(apperror.CodeParamFormat, "支付渠道配置格式错误", err)
	}
	return cfg, nil
}

func epaySign(params map[string]string, key string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, k := range keys {
		v := params[k]
		if v == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteByte('&')
		}
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(v)
	}
	b.WriteString("&key=")
	b.WriteString(key)
	sum := md5.Sum([]byte(b.String()))
	return hex.EncodeToString(sum[:])
}
