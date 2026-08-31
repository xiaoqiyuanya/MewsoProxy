package subscribe

import (
	"context"
	"encoding/base64"
	"fmt"

	"mewsoproxy/server/config"
	"mewsoproxy/server/model"
)

type Service struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Service {
	return &Service{cfg: cfg}
}

type Result struct {
	Token string `json:"token"`
	URL   string `json:"url"`
}

func (s *Service) Generate(ctx context.Context, u *model.User) Result {
	url := fmt.Sprintf("%s/subscribe?token=%s", s.cfg.App.SubscribeURL, u.Token)
	return Result{Token: u.Token, URL: url}
}

func (s *Service) EncodeURL(raw string) string {
	return base64.StdEncoding.EncodeToString([]byte(raw))
}
