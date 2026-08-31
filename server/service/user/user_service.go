package user

import (
	"context"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"mewsoproxy/server/config"
	"mewsoproxy/server/dto"
	"mewsoproxy/server/model"
	"mewsoproxy/server/pkg/apperror"
	"mewsoproxy/server/pkg/token"
)

type Service struct {
	db  *gorm.DB
	cfg *config.Config
}

func New(db *gorm.DB, cfg *config.Config) *Service {
	return &Service{db: db, cfg: cfg}
}

const BcryptCost = 12

func (s *Service) Register(ctx context.Context, req dto.RegisterReq) (*model.User, error) {
	if !s.cfg.App.RegisterEnabled {
		return nil, apperror.New(apperror.CodeNoPermission, "注册已关闭")
	}
	var cnt int64
	if err := s.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", req.Email).Count(&cnt).Error; err != nil {
		return nil, apperror.Wrap(apperror.CodeDBError, "数据库错误", err)
	}
	if cnt > 0 {
		return nil, apperror.New(apperror.CodeUserExists, "邮箱已被注册")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), BcryptCost)
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeDBError, "密码处理失败", err)
	}
	inviteUID := uint(0)
	if req.InviteCode != "" {
		var code model.InviteCode
		if err := s.db.WithContext(ctx).Where("code = ?", req.InviteCode).First(&code).Error; err == nil {
			inviteUID = code.UserID
		}
	}
	u := &model.User{
		Email:            req.Email,
		Password:         string(hash),
		UUID:             newUUID(),
		Token:            newToken32(),
		Balance:          0,
		CommissionType:   0,
		CommissionRate:   intp(s.defaultCommissionRate()),
		TransferEnable:   0,
		GroupID:          intp(s.cfg.App.DefaultGroupID),
		PlanID:           nil,
		Banned:           false,
		IsAdmin:          false,
		CreatedAt:        nowUnix(),
		UpdatedAt:        nowUnix(),
		ExpiredAt:        0,
	}
	if inviteUID > 0 {
		u.InviteUserID = uintp(inviteUID)
	}
	if err := s.db.WithContext(ctx).Create(u).Error; err != nil {
		return nil, apperror.Wrap(apperror.CodeDBError, "创建用户失败", err)
	}
	return u, nil
}

func (s *Service) Login(ctx context.Context, req dto.LoginReq) (*model.User, error) {
	var u model.User
	if err := s.db.WithContext(ctx).Where("email = ?", req.Email).First(&u).Error; err != nil {
		return nil, apperror.New(apperror.CodeUserNotFound, "账号不存在")
	}
	if u.Banned {
		return nil, apperror.New(apperror.CodeNoPermission, "账号已被封禁")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, apperror.New(apperror.CodePasswordWrong, "密码错误")
	}
	now := nowUnix()
	u.LastLoginAt = &now
	if err := s.db.WithContext(ctx).Model(&u).Update("last_login_at", now).Error; err != nil {
		return nil, apperror.Wrap(apperror.CodeDBError, "登录记录失败", err)
	}
	return &u, nil
}

func (s *Service) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var u model.User
	if err := s.db.WithContext(ctx).First(&u, id).Error; err != nil {
		return nil, apperror.New(apperror.CodeUserNotFound, "用户不存在")
	}
	if u.Banned {
		return nil, apperror.New(apperror.CodeNoPermission, "账号已被封禁")
	}
	return &u, nil
}

func (s *Service) GetByToken(ctx context.Context, token string) (*model.User, error) {
	if token == "" {
		return nil, apperror.New(apperror.CodeParamMissing, "订阅令牌缺失")
	}
	var u model.User
	if err := s.db.WithContext(ctx).Where("token = ?", token).First(&u).Error; err != nil {
		return nil, apperror.New(apperror.CodeUserNotFound, "订阅令牌无效")
	}
	if u.Banned {
		return nil, apperror.New(apperror.CodeNoPermission, "账号已被封禁")
	}
	return &u, nil
}

func (s *Service) defaultCommissionRate() int {
	return 0
}

func nowUnix() int64 {
	return time.Now().UTC().Unix()
}

func newUUID() string {
	s, _ := token.RandomString(16)
	return s
}

func newToken32() string {
	s, _ := token.RandomString(16)
	return s
}

func intp(i int) *int    { return &i }
func uintp(u uint) *uint { return &u }

func ToDTO(u *model.User) dto.UserDTO {
	return dto.UserDTO{
		ID:               u.ID,
		Email:            u.Email,
		Balance:          u.Balance,
		CommissionBalance: u.CommissionBalance,
		IsAdmin:          u.IsAdmin,
		IsStaff:          u.IsStaff,
		Banned:           u.Banned,
		PlanID:           u.PlanID,
		GroupID:          u.GroupID,
		ExpiredAt:        dto.FromUnix(u.ExpiredAt),
		Token:            u.Token,
		UUID:             u.UUID,
		TransferEnable:   u.TransferEnable,
		UsedTraffic:      u.U + u.D,
		CreatedAt:        dto.FromUnix(u.CreatedAt),
	}
}
