package admin

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"mewsoproxy/server/dto"
	"mewsoproxy/server/model"
	"mewsoproxy/server/pkg/apperror"
	"mewsoproxy/server/pkg/crypto"
	"mewsoproxy/server/pkg/sshclient"
)

type installNode struct {
	id       int
	nodeType string
	name     string
	serverPort int
	cipher   string
	network  string
	tls      bool
	sni      string
	upMbps   int
	downMbps int
	insecure bool
}

type installTask struct {
	id          string
	mu          sync.Mutex
	logs        []string
	subscribers map[chan string]struct{}
	done        bool
	errMsg      string
	doneCh      chan struct{}
}

func (t *installTask) push(line string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.logs = append(t.logs, line)
	for ch := range t.subscribers {
		select {
		case ch <- line:
		default:
		}
	}
}

func (t *installTask) finish(err error) {
	t.mu.Lock()
	t.done = true
	if err != nil {
		t.errMsg = err.Error()
		line := "##INSTALL_FAILED: " + err.Error()
		t.logs = append(t.logs, line)
		for ch := range t.subscribers {
			select {
			case ch <- line:
			default:
			}
		}
	} else {
		t.logs = append(t.logs, "##INSTALL_DONE##")
		for ch := range t.subscribers {
			select {
			case ch <- "##INSTALL_DONE##":
			default:
			}
		}
	}
	close(t.doneCh)
	t.mu.Unlock()
}

func (t *installTask) unsubscribe(ch chan string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.subscribers, ch)
}

type taskStore struct {
	mu    sync.Mutex
	tasks map[string]*installTask
}

var installTasks = &taskStore{tasks: map[string]*installTask{}}

func (ts *taskStore) add(t *installTask) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.tasks[t.id] = t
}

func (ts *taskStore) get(id string) (*installTask, bool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	t, ok := ts.tasks[id]
	return t, ok
}

func (s *Service) InstallNode(ctx context.Context, req dto.AdminNodeInstallReq) (string, error) {
	node, err := s.loadInstallNode(ctx, req.Type, req.ID)
	if err != nil {
		return "", err
	}
	taskID := uuid.New().String()
	t := &installTask{id: taskID, subscribers: map[chan string]struct{}{}, doneCh: make(chan struct{})}
	installTasks.add(t)

	sshRec, err := s.upsertSSHRecord(ctx, req)
	if err != nil {
		return "", err
	}
	s.updateInstallStatus(ctx, sshRec, model.NodeInstallStatusRunning, "")

	payload := installPayload{node: node, req: req, sshRec: sshRec, uuid: s.pickNodeUUID(ctx, node)}
	go s.runInstall(taskID, payload)
	return taskID, nil
}

type installPayload struct {
	node   installNode
	req    dto.AdminNodeInstallReq
	sshRec *model.ServerNodeSSH
	uuid   string
}

func (s *Service) SubscribeLog(taskID string) (<-chan string, func(), error) {
	t, ok := installTasks.get(taskID)
	if !ok {
		return nil, nil, apperror.New(apperror.CodeResourceNotFnd, "安装任务不存在")
	}
	ch := make(chan string, 256)
	t.mu.Lock()
	if t.done {
		for _, l := range t.logs {
			ch <- l
		}
		close(ch)
		t.mu.Unlock()
		return ch, func() {}, nil
	}
	for _, l := range t.logs {
		ch <- l
	}
	t.subscribers[ch] = struct{}{}
	t.mu.Unlock()
	return ch, func() { t.unsubscribe(ch) }, nil
}

func (s *Service) runInstall(taskID string, p installPayload) {
	t, ok := installTasks.get(taskID)
	if !ok {
		return
	}
	t.push(fmt.Sprintf("开始连接节点 %s:%d ...", p.sshRec.SSHHost, p.sshRec.SSHPort))
	password, key, err := s.decryptSSH(p.sshRec)
	if err != nil {
		t.finish(err)
		s.updateInstallStatus(context.Background(), p.sshRec, model.NodeInstallStatusFailed, err.Error())
		return
	}
	cli, err := sshclient.New(sshclient.Config{
		Host:       p.sshRec.SSHHost,
		Port:       p.sshRec.SSHPort,
		User:       p.sshRec.SSHUser,
		Password:   password,
		PrivateKey: key,
		Timeout:    30 * time.Second,
	})
	if err != nil {
		t.finish(err)
		s.updateInstallStatus(context.Background(), p.sshRec, model.NodeInstallStatusFailed, err.Error())
		return
	}
	defer cli.Close()

	script := buildInstallScript(p.node, p.uuid)
	cmd := fmt.Sprintf("cat > /tmp/mewso_install.sh <<'MEWSOEOF'\n%s\nMEWSOEOF\nbash /tmp/mewso_install.sh", script)
	err = cli.RunWithLog(context.Background(), cmd, t.push)
	if err != nil {
		t.finish(err)
		s.updateInstallStatus(context.Background(), p.sshRec, model.NodeInstallStatusFailed, err.Error())
		return
	}
	t.finish(nil)
	s.updateInstallStatus(context.Background(), p.sshRec, model.NodeInstallStatusSuccess, "安装完成")
}

func (s *Service) loadInstallNode(ctx context.Context, nodeType string, id int) (installNode, error) {
	var n installNode
	n.nodeType = nodeType
	n.id = id
	switch nodeType {
	case "trojan":
		var m model.ServerTrojan
		if err := s.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
			return n, apperror.New(apperror.CodeResourceNotFnd, "节点不存在")
		}
		n.name, n.serverPort, n.tls, n.sni = m.Name, m.ServerPort, true, strOr(m.ServerName)
	case "vmess":
		var m model.ServerVmess
		if err := s.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
			return n, apperror.New(apperror.CodeResourceNotFnd, "节点不存在")
		}
		n.name, n.serverPort, n.tls, n.network, n.sni = m.Name, m.ServerPort, m.TLS == 1, m.Network, m.Host
	case "shadowsocks", "ss":
		var m model.ServerShadowsocks
		if err := s.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
			return n, apperror.New(apperror.CodeResourceNotFnd, "节点不存在")
		}
		n.name, n.serverPort, n.cipher = m.Name, m.ServerPort, m.Cipher
	case "hysteria":
		var m model.ServerHysteria
		if err := s.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
			return n, apperror.New(apperror.CodeResourceNotFnd, "节点不存在")
		}
		n.name, n.serverPort, n.sni, n.insecure = m.Name, m.ServerPort, strOr(m.ServerName), m.Insecure
		n.upMbps, n.downMbps = m.UpMbps, m.DownMbps
	default:
		return n, apperror.New(apperror.CodeParamFormat, "不支持的节点类型")
	}
	return n, nil
}

func (s *Service) upsertSSHRecord(ctx context.Context, req dto.AdminNodeInstallReq) (*model.ServerNodeSSH, error) {
	var rec model.ServerNodeSSH
	if err := s.db.WithContext(ctx).Where("node_type = ? AND node_id = ?", req.Type, req.ID).First(&rec).Error; err != nil {
		now := time.Now().UTC().Unix()
		rec = model.ServerNodeSSH{NodeType: req.Type, NodeID: req.ID, SSHPort: 22, CreatedAt: now, UpdatedAt: now}
	}
	if req.SSHHost != "" {
		rec.SSHHost = req.SSHHost
	}
	if req.SSHPort > 0 {
		rec.SSHPort = req.SSHPort
	}
	if req.SSHUser != "" {
		rec.SSHUser = req.SSHUser
	}
	if req.SSHPassword != "" {
		enc, err := crypto.Encrypt(s.cfg.App.SSHEncryptKey, req.SSHPassword)
		if err != nil {
			return nil, apperror.Wrap(apperror.CodeParamFormat, "SSH 密码加密失败", err)
		}
		rec.SSHPassword = &enc
	}
	if req.SSHPrivateKey != "" {
		enc, err := crypto.Encrypt(s.cfg.App.SSHEncryptKey, req.SSHPrivateKey)
		if err != nil {
			return nil, apperror.Wrap(apperror.CodeParamFormat, "SSH 私钥加密失败", err)
		}
		rec.SSHPrivateKey = &enc
	}
	if rec.SSHHost == "" || rec.SSHUser == "" {
		return nil, apperror.New(apperror.CodeParamMissing, "SSH 地址或用户名缺失")
	}
	now := time.Now().UTC().Unix()
	rec.UpdatedAt = now
	if rec.ID == 0 {
		if err := s.db.WithContext(ctx).Create(&rec).Error; err != nil {
			return nil, apperror.Wrap(apperror.CodeDBError, "保存节点 SSH 配置失败", err)
		}
	} else {
		if err := s.db.WithContext(ctx).Save(&rec).Error; err != nil {
			return nil, apperror.Wrap(apperror.CodeDBError, "保存节点 SSH 配置失败", err)
		}
	}
	return &rec, nil
}

func (s *Service) updateInstallStatus(ctx context.Context, rec *model.ServerNodeSSH, status int8, log string) {
	now := time.Now().UTC().Unix()
	updates := model.ServerNodeSSH{InstallStatus: status, UpdatedAt: now}
	if log != "" {
		updates.LastLog = &log
	}
	s.db.WithContext(ctx).Model(&model.ServerNodeSSH{}).Where("id = ?", rec.ID).Updates(updates)
}

func (s *Service) pickNodeUUID(ctx context.Context, node installNode) string {
	var u model.User
	if err := s.db.WithContext(ctx).Where("banned = 0").Order("id asc").First(&u).Error; err == nil && u.UUID != "" {
		return u.UUID
	}
	ru, _ := uuid.NewRandom()
	return ru.String()
}

func (s *Service) decryptSSH(rec *model.ServerNodeSSH) (password, key string, err error) {
	if rec.SSHPassword != nil {
		if password, err = crypto.Decrypt(s.cfg.App.SSHEncryptKey, *rec.SSHPassword); err != nil {
			return "", "", err
		}
	}
	if rec.SSHPrivateKey != nil {
		if key, err = crypto.Decrypt(s.cfg.App.SSHEncryptKey, *rec.SSHPrivateKey); err != nil {
			return "", "", err
		}
	}
	return password, key, nil
}

func strOr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
