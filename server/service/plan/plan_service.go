package plan

import (
	"context"

	"gorm.io/gorm"

	"mewsoproxy/server/dto"
	"mewsoproxy/server/model"
	"mewsoproxy/server/pkg/apperror"
)

type Service struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) List(ctx context.Context) ([]dto.PlanDetailDTO, error) {
	var plans []model.Plan
	if err := s.db.WithContext(ctx).Where("`show` = ?", true).Order("sort asc").Find(&plans).Error; err != nil {
		return nil, apperror.Wrap(apperror.CodeDBError, "套餐列表获取失败", err)
	}
	out := make([]dto.PlanDetailDTO, 0, len(plans))
	for i := range plans {
		out = append(out, ToDTO(&plans[i]))
	}
	return out, nil
}

func (s *Service) Get(ctx context.Context, id int) (*model.Plan, error) {
	var p model.Plan
	if err := s.db.WithContext(ctx).First(&p, id).Error; err != nil {
		return nil, apperror.New(apperror.CodeResourceNotFnd, "套餐不存在")
	}
	return &p, nil
}

func ToDTO(p *model.Plan) dto.PlanDetailDTO {
	return dto.PlanDetailDTO{
		ID:             p.ID,
		Name:           p.Name,
		GroupID:        p.GroupID,
		TransferEnable: p.TransferEnable,
		SpeedLimit:     p.SpeedLimit,
		Show:           p.Show,
		Sort:           p.Sort,
		Renew:          p.Renew,
		Content:        p.Content,
		MonthPrice:     p.MonthPrice,
		QuarterPrice:   p.QuarterPrice,
		HalfYearPrice:  p.HalfYearPrice,
		YearPrice:      p.YearPrice,
		TwoYearPrice:   p.TwoYearPrice,
		ThreeYearPrice: p.ThreeYearPrice,
		OnetimePrice:   p.OnetimePrice,
	}
}
