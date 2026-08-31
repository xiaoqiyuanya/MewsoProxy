package subscribe

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"mewsoproxy/server/config"
	"mewsoproxy/server/model"
	"mewsoproxy/server/pkg/apperror"
)

type Service struct {
	cfg *config.Config
	db  *gorm.DB
}

func New(cfg *config.Config, db *gorm.DB) *Service {
	return &Service{cfg: cfg, db: db}
}

type Result struct {
	Token string `json:"token"`
	URL   string `json:"url"`
}

type ContentType string

const (
	ContentV2ray ContentType = "v2ray"
	ContentClash ContentType = "clash"
)

func (s *Service) Generate(ctx context.Context, u *model.User) Result {
	url := fmt.Sprintf("%s/subscribe?token=%s", s.cfg.App.SubscribeURL, u.Token)
	return Result{Token: u.Token, URL: url}
}

func (s *Service) BuildSubscription(ctx context.Context, u *model.User, ctype ContentType) (string, error) {
	if u.GroupID == nil {
		return "", apperror.New(apperror.CodeResourceNotFnd, "该用户未分配节点分组")
	}
	nodes, err := s.collectNodes(ctx, *u.GroupID)
	if err != nil {
		return "", err
	}
	if len(nodes) == 0 {
		return "", apperror.New(apperror.CodeResourceNotFnd, "当前分组下暂无可用节点")
	}
	switch ctype {
	case ContentClash:
		return buildClash(nodes, u.UUID), nil
	default:
		return buildV2ray(nodes, u.UUID), nil
	}
}

func EncodeURL(raw string) string {
	return base64Encode(raw)
}
