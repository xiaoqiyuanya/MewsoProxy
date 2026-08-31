package admin

import (
	"context"
	"time"

	"mewsoproxy/server/dto"
	"mewsoproxy/server/model"
	"mewsoproxy/server/pkg/apperror"
)

func (s *Service) ListGroups(ctx context.Context) ([]model.ServerGroup, error) {
	var list []model.ServerGroup
	if err := s.db.WithContext(ctx).Order("id asc").Find(&list).Error; err != nil {
		return nil, apperror.Wrap(apperror.CodeDBError, "节点组查询失败", err)
	}
	return list, nil
}

func (s *Service) SaveGroup(ctx context.Context, req dto.AdminServerGroupSaveReq) (int, error) {
	now := time.Now().UTC().Unix()
	if req.ID > 0 {
		if err := s.db.WithContext(ctx).Model(&model.ServerGroup{}).Where("id = ?", req.ID).
			Update("name", req.Name).Update("updated_at", now).Error; err != nil {
			return 0, apperror.Wrap(apperror.CodeDBError, "保存节点组失败", err)
		}
		return req.ID, nil
	}
	g := model.ServerGroup{Name: req.Name, CreatedAt: now, UpdatedAt: now}
	if err := s.db.WithContext(ctx).Create(&g).Error; err != nil {
		return 0, apperror.Wrap(apperror.CodeDBError, "创建节点组失败", err)
	}
	return g.ID, nil
}

func (s *Service) DropGroup(ctx context.Context, id int) error {
	res := s.db.WithContext(ctx).Where("id = ?", id).Delete(&model.ServerGroup{})
	if res.Error != nil {
		return apperror.Wrap(apperror.CodeDBError, "删除节点组失败", res.Error)
	}
	if res.RowsAffected == 0 {
		return apperror.New(apperror.CodeResourceNotFnd, "节点组不存在")
	}
	return nil
}

func (s *Service) ListNodes(ctx context.Context, nodeType string) (interface{}, error) {
	switch nodeType {
	case "trojan":
		var list []model.ServerTrojan
		if err := s.db.WithContext(ctx).Order("sort asc").Find(&list).Error; err != nil {
			return nil, apperror.Wrap(apperror.CodeDBError, "节点查询失败", err)
		}
		return list, nil
	case "vmess":
		var list []model.ServerVmess
		if err := s.db.WithContext(ctx).Order("sort asc").Find(&list).Error; err != nil {
			return nil, apperror.Wrap(apperror.CodeDBError, "节点查询失败", err)
		}
		return list, nil
	case "shadowsocks", "ss":
		var list []model.ServerShadowsocks
		if err := s.db.WithContext(ctx).Order("sort asc").Find(&list).Error; err != nil {
			return nil, apperror.Wrap(apperror.CodeDBError, "节点查询失败", err)
		}
		return list, nil
	case "hysteria":
		var list []model.ServerHysteria
		if err := s.db.WithContext(ctx).Order("sort asc").Find(&list).Error; err != nil {
			return nil, apperror.Wrap(apperror.CodeDBError, "节点查询失败", err)
		}
		return list, nil
	default:
		return nil, apperror.New(apperror.CodeParamFormat, "不支持的节点类型")
	}
}

func (s *Service) SaveNode(ctx context.Context, req dto.AdminServerNodeSaveReq) (int, error) {
	now := time.Now().UTC().Unix()
	switch req.Type {
	case "trojan":
		m := toTrojan(req)
		if req.ID > 0 {
			m.UpdatedAt = now
			if err := s.db.WithContext(ctx).Model(&model.ServerTrojan{}).Where("id = ?", req.ID).Updates(&m).Error; err != nil {
				return 0, apperror.Wrap(apperror.CodeDBError, "保存节点失败", err)
			}
			return req.ID, nil
		}
		m.CreatedAt = now
		if err := s.db.WithContext(ctx).Create(&m).Error; err != nil {
			return 0, apperror.Wrap(apperror.CodeDBError, "创建节点失败", err)
		}
		return m.ID, nil
	case "vmess":
		m := toVmess(req)
		if req.ID > 0 {
			m.UpdatedAt = now
			if err := s.db.WithContext(ctx).Model(&model.ServerVmess{}).Where("id = ?", req.ID).Updates(&m).Error; err != nil {
				return 0, apperror.Wrap(apperror.CodeDBError, "保存节点失败", err)
			}
			return req.ID, nil
		}
		m.CreatedAt = now
		if err := s.db.WithContext(ctx).Create(&m).Error; err != nil {
			return 0, apperror.Wrap(apperror.CodeDBError, "创建节点失败", err)
		}
		return m.ID, nil
	case "shadowsocks", "ss":
		m := toSS(req)
		if req.ID > 0 {
			m.UpdatedAt = now
			if err := s.db.WithContext(ctx).Model(&model.ServerShadowsocks{}).Where("id = ?", req.ID).Updates(&m).Error; err != nil {
				return 0, apperror.Wrap(apperror.CodeDBError, "保存节点失败", err)
			}
			return req.ID, nil
		}
		m.CreatedAt = now
		if err := s.db.WithContext(ctx).Create(&m).Error; err != nil {
			return 0, apperror.Wrap(apperror.CodeDBError, "创建节点失败", err)
		}
		return m.ID, nil
	case "hysteria":
		m := toHysteria(req)
		if req.ID > 0 {
			m.UpdatedAt = now
			if err := s.db.WithContext(ctx).Model(&model.ServerHysteria{}).Where("id = ?", req.ID).Updates(&m).Error; err != nil {
				return 0, apperror.Wrap(apperror.CodeDBError, "保存节点失败", err)
			}
			return req.ID, nil
		}
		m.CreatedAt = now
		if err := s.db.WithContext(ctx).Create(&m).Error; err != nil {
			return 0, apperror.Wrap(apperror.CodeDBError, "创建节点失败", err)
		}
		return m.ID, nil
	default:
		return 0, apperror.New(apperror.CodeParamFormat, "不支持的节点类型")
	}
}

func (s *Service) DropNode(ctx context.Context, nodeType string, id int) error {
	var err error
	switch nodeType {
	case "trojan":
		err = s.db.WithContext(ctx).Where("id = ?", id).Delete(&model.ServerTrojan{}).Error
	case "vmess":
		err = s.db.WithContext(ctx).Where("id = ?", id).Delete(&model.ServerVmess{}).Error
	case "shadowsocks", "ss":
		err = s.db.WithContext(ctx).Where("id = ?", id).Delete(&model.ServerShadowsocks{}).Error
	case "hysteria":
		err = s.db.WithContext(ctx).Where("id = ?", id).Delete(&model.ServerHysteria{}).Error
	default:
		return apperror.New(apperror.CodeParamFormat, "不支持的节点类型")
	}
	if err != nil {
		return apperror.Wrap(apperror.CodeDBError, "删除节点失败", err)
	}
	return nil
}

func toTrojan(r dto.AdminServerNodeSaveReq) model.ServerTrojan {
	return model.ServerTrojan{
		GroupID: r.GroupID, RouteID: r.RouteID, ParentID: r.ParentID, Tags: r.Tags,
		Name: r.Name, Rate: defaultStr(r.Rate, "1"), Host: r.Host, Port: r.Port,
		ServerPort: r.ServerPort, Show: r.Show, Sort: r.Sort,
		AllowInsecure: boolOr(r.AllowInsecure, false), ServerName: r.ServerName,
	}
}

func toVmess(r dto.AdminServerNodeSaveReq) model.ServerVmess {
	return model.ServerVmess{
		GroupID: r.GroupID, RouteID: r.RouteID, ParentID: r.ParentID, Tags: r.Tags,
		Name: r.Name, Host: r.Host, Port: r.Port, ServerPort: r.ServerPort,
		TLS: int8Or(r.TLS, 0), Rate: defaultStr(r.Rate, "1"), Network: strPtrOr(r.Network, "tcp"),
		Show: r.Show, Sort: r.Sort,
	}
}

func toSS(r dto.AdminServerNodeSaveReq) model.ServerShadowsocks {
	return model.ServerShadowsocks{
		GroupID: r.GroupID, RouteID: r.RouteID, ParentID: r.ParentID, Tags: r.Tags,
		Name: r.Name, Host: r.Host, Port: r.Port, ServerPort: r.ServerPort,
		Cipher: strPtrOr(r.Cipher, "aes-256-gcm"), Rate: defaultStr(r.Rate, "1"),
		Show: r.Show, Sort: r.Sort,
	}
}

func toHysteria(r dto.AdminServerNodeSaveReq) model.ServerHysteria {
	return model.ServerHysteria{
		GroupID: r.GroupID, RouteID: r.RouteID, ParentID: r.ParentID, Tags: r.Tags,
		Name: r.Name, Host: r.Host, Port: r.Port, ServerPort: r.ServerPort,
		Rate: defaultStr(r.Rate, "1"), Show: r.Show, Sort: r.Sort,
		UpMbps: intOr(r.UpMbps, 0), DownMbps: intOr(r.DownMbps, 0),
		ServerName: r.ServerName, Insecure: boolOr(r.Insecure, false),
	}
}

func defaultStr(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

func strPtrOr(p *string, def string) string {
	if p == nil || *p == "" {
		return def
	}
	return *p
}

func boolOr(p *bool, def bool) bool {
	if p == nil {
		return def
	}
	return *p
}

func int8Or(p *int8, def int8) int8 {
	if p == nil {
		return def
	}
	return *p
}

func intOr(p *int, def int) int {
	if p == nil {
		return def
	}
	return *p
}
