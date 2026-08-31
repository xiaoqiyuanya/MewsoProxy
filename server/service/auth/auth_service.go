package auth

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"mewsoproxy/server/config"
	"mewsoproxy/server/dto"
	"mewsoproxy/server/model"
	"mewsoproxy/server/pkg/apperror"
	redisc "mewsoproxy/server/pkg/redis"
	"mewsoproxy/server/pkg/token"
	usersvc "mewsoproxy/server/service/user"
)

type Service struct {
	cfg     *config.Config
	rds     *redisc.Client
	userSvc *usersvc.Service
}

func New(cfg *config.Config, rds *redisc.Client, userSvc *usersvc.Service) *Service {
	return &Service{cfg: cfg, rds: rds, userSvc: userSvc}
}

const maxRefreshTokens = 5

func (s *Service) Issue(ctx context.Context, u *model.User) (dto.TokenDTO, error) {
	access, err := token.GenerateAccessToken(u.ID, u.IsAdmin, s.cfg.JWT.AccessSecret, s.parseDuration(s.cfg.JWT.AccessTTL))
	if err != nil {
		return dto.TokenDTO{}, apperror.Wrap(apperror.CodeDBError, "签发令牌失败", err)
	}
	refreshToken, err := token.RandomString(32)
	if err != nil {
		return dto.TokenDTO{}, apperror.Wrap(apperror.CodeDBError, "签发令牌失败", err)
	}
	tokenID, _ := token.RandomString(8)
	refreshKey := redisc.RefreshTokenKey(u.ID, tokenID)
	indexKey := redisc.RefreshIndexKey(refreshToken)
	refreshTTL := s.parseDuration(s.cfg.JWT.RefreshTTL)

	if err := s.rds.Set(ctx, refreshKey, refreshToken, refreshTTL); err != nil {
		return dto.TokenDTO{}, apperror.Wrap(apperror.CodeDBError, "令牌存储失败", err)
	}
	if err := s.rds.Set(ctx, indexKey, fmt.Sprintf("%d:%s", u.ID, tokenID), refreshTTL); err != nil {
		return dto.TokenDTO{}, apperror.Wrap(apperror.CodeDBError, "令牌存储失败", err)
	}
	if err := s.pushTokenID(ctx, u.ID, tokenID, refreshTTL); err != nil {
		return dto.TokenDTO{}, err
	}
	dur := int64(access.ExpiresAt.Sub(time.Now()).Seconds())
	return dto.TokenDTO{AccessToken: access.Token, RefreshToken: refreshToken, ExpiresIn: dur}, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (dto.TokenDTO, error) {
	indexKey := redisc.RefreshIndexKey(refreshToken)
	val, err := s.rds.Get(ctx, indexKey)
	if err != nil || val == "" {
		return dto.TokenDTO{}, apperror.New(apperror.CodeTokenExpired, "登录已过期，请重新登录")
	}
	parts := strings.SplitN(val, ":", 2)
	if len(parts) != 2 {
		return dto.TokenDTO{}, apperror.New(apperror.CodeTokenExpired, "登录已过期，请重新登录")
	}
	uid64, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return dto.TokenDTO{}, apperror.New(apperror.CodeTokenExpired, "登录已过期，请重新登录")
	}
	uid := uint(uid64)
	tokenID := parts[1]
	refreshKey := redisc.RefreshTokenKey(uid, tokenID)
	actual, err := s.rds.Get(ctx, refreshKey)
	if err != nil || actual != refreshToken {
		return dto.TokenDTO{}, apperror.New(apperror.CodeTokenExpired, "登录已过期，请重新登录")
	}
	u, err := s.userSvc.GetByID(ctx, uid)
	if err != nil {
		return dto.TokenDTO{}, err
	}
	if err := s.rotateToken(ctx, uid, tokenID, indexKey); err != nil {
		return dto.TokenDTO{}, apperror.Wrap(apperror.CodeDBError, "令牌轮换失败", err)
	}
	return s.Issue(ctx, u)
}

func (s *Service) Logout(ctx context.Context, refreshToken string, accTokenExp time.Time) error {
	if refreshToken != "" {
		val, err := s.rds.Get(ctx, redisc.RefreshIndexKey(refreshToken))
		if err == nil && val != "" {
			parts := strings.SplitN(val, ":", 2)
			if len(parts) == 2 {
				uid64, _ := strconv.ParseUint(parts[0], 10, 64)
				tokenID := parts[1]
				_ = s.rds.Del(ctx, redisc.RefreshTokenKey(uint(uid64), tokenID), redisc.RefreshIndexKey(refreshToken), redisc.UserRefreshKey(uint(uid64)))
			}
		}
	}
	return nil
}

func (s *Service) Blacklist(ctx context.Context, jti string, exp time.Time) error {
	ttl := time.Until(exp)
	if ttl <= 0 {
		ttl = time.Minute
	}
	return s.rds.Set(ctx, redisc.TokenBlacklistKey(jti), "1", ttl)
}

func (s *Service) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	return s.rds.Exists(ctx, redisc.TokenBlacklistKey(jti))
}

func (s *Service) rotateToken(ctx context.Context, uid uint, oldID, oldIndex string) error {
	_ = s.rds.Del(ctx, redisc.RefreshTokenKey(uid, oldID), oldIndex)
	if items, err := s.rds.LRange(ctx, redisc.UserRefreshKey(uid), 0, -1); err == nil {
		var kept []interface{}
		for _, it := range items {
			if it != oldID {
				kept = append(kept, it)
			}
		}
		_ = s.rds.Del(ctx, redisc.UserRefreshKey(uid))
		if len(kept) > 0 {
			_ = s.rds.RPush(ctx, redisc.UserRefreshKey(uid), kept...)
		}
	}
	return nil
}

func (s *Service) pushTokenID(ctx context.Context, uid uint, tokenID string, ttl time.Duration) error {
	if items, err := s.rds.LRange(ctx, redisc.UserRefreshKey(uid), 0, -1); err == nil {
		if len(items) >= maxRefreshTokens {
			_ = s.rds.Del(ctx, redisc.UserRefreshKey(uid))
		}
	}
	return nil
}

func (s *Service) parseDuration(v string) time.Duration {
	d, err := time.ParseDuration(v)
	if err != nil {
		return 30 * time.Minute
	}
	return d
}
