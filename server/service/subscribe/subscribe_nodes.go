package subscribe

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"mewsoproxy/server/model"
	"mewsoproxy/server/pkg/apperror"
)

type nodeInfo struct {
	typ           string
	name          string
	host          string
	port          string
	serverPort    int
	rate          string
	tls           bool
	network       string
	sni           string
	cipher        string
	allowInsecure bool
	upMbps        int
	downMbps      int
	insecure      bool
}

func (s *Service) collectNodes(ctx context.Context, groupID int) ([]nodeInfo, error) {
	var nodes []nodeInfo

	var trojans []model.ServerTrojan
	if err := s.db.WithContext(ctx).Find(&trojans).Error; err != nil {
		return nil, apperror.Wrap(apperror.CodeDBError, "节点查询失败", err)
	}
	for _, n := range trojans {
		if !n.Show || !groupContains(n.GroupID, groupID) {
			continue
		}
		nodes = append(nodes, nodeInfo{
			typ: "trojan", name: n.Name, host: n.Host, port: n.Port, serverPort: n.ServerPort,
			rate: n.Rate, tls: true, network: "tcp", sni: valOr(n.ServerName),
			allowInsecure: n.AllowInsecure,
		})
	}

	var vmess []model.ServerVmess
	if err := s.db.WithContext(ctx).Find(&vmess).Error; err != nil {
		return nil, apperror.Wrap(apperror.CodeDBError, "节点查询失败", err)
	}
	for _, n := range vmess {
		if !n.Show || !groupContains(n.GroupID, groupID) {
			continue
		}
		nodes = append(nodes, nodeInfo{
			typ: "vmess", name: n.Name, host: n.Host, port: n.Port, serverPort: n.ServerPort,
			rate: n.Rate, tls: n.TLS == 1, network: n.Network, sni: n.Host,
		})
	}

	var ss []model.ServerShadowsocks
	if err := s.db.WithContext(ctx).Find(&ss).Error; err != nil {
		return nil, apperror.Wrap(apperror.CodeDBError, "节点查询失败", err)
	}
	for _, n := range ss {
		if !n.Show || !groupContains(n.GroupID, groupID) {
			continue
		}
		nodes = append(nodes, nodeInfo{
			typ: "shadowsocks", name: n.Name, host: n.Host, port: n.Port, serverPort: n.ServerPort,
			rate: n.Rate, cipher: n.Cipher,
		})
	}

	var hys []model.ServerHysteria
	if err := s.db.WithContext(ctx).Find(&hys).Error; err != nil {
		return nil, apperror.Wrap(apperror.CodeDBError, "节点查询失败", err)
	}
	for _, n := range hys {
		if !n.Show || !groupContains(n.GroupID, groupID) {
			continue
		}
		nodes = append(nodes, nodeInfo{
			typ: "hysteria", name: n.Name, host: n.Host, port: n.Port, serverPort: n.ServerPort,
			rate: n.Rate, sni: valOr(n.ServerName), upMbps: n.UpMbps, downMbps: n.DownMbps,
			insecure: n.Insecure,
		})
	}
	return nodes, nil
}

func groupContains(list string, gid int) bool {
	if list == "" {
		return false
	}
	target := strconv.Itoa(gid)
	for _, part := range strings.Split(list, ",") {
		if strings.TrimSpace(part) == target {
			return true
		}
	}
	return false
}

func buildV2ray(nodes []nodeInfo, uuid string) string {
	var lines []string
	for _, n := range nodes {
		if link := nodeShareLink(n, uuid); link != "" {
			lines = append(lines, link)
		}
	}
	return base64.StdEncoding.EncodeToString([]byte(strings.Join(lines, "\n")))
}

func nodeShareLink(n nodeInfo, uuid string) string {
	name := url.QueryEscape(n.name)
	switch n.typ {
	case "trojan":
		q := url.Values{}
		q.Set("security", "tls")
		if n.sni != "" {
			q.Set("sni", n.sni)
		}
		q.Set("type", "tcp")
		if n.allowInsecure {
			q.Set("allowInsecure", "1")
		}
		return fmt.Sprintf("trojan://%s@%s:%s?%s#%s", uuid, n.host, n.port, q.Encode(), name)
	case "vmess":
		obj := map[string]interface{}{
			"v": 2, "ps": n.name, "add": n.host, "port": n.port, "id": uuid,
			"aid": 0, "net": n.network, "type": "none",
		}
		if n.tls {
			obj["tls"] = "tls"
			obj["sni"] = n.sni
		}
		return "vmess://" + base64.RawURLEncoding.EncodeToString([]byte(vmessJSON(obj)))
	case "shadowsocks":
		return fmt.Sprintf("ss://%s@%s:%s#%s", ssPayload(n.cipher, uuid), n.host, n.port, name)
	case "hysteria":
		q := url.Values{}
		q.Set("auth", uuid)
		if n.sni != "" {
			q.Set("peer", n.sni)
			q.Set("sni", n.sni)
		}
		if n.upMbps > 0 {
			q.Set("upmbps", strconv.Itoa(n.upMbps))
		}
		if n.downMbps > 0 {
			q.Set("downmbps", strconv.Itoa(n.downMbps))
		}
		if n.insecure {
			q.Set("insecure", "1")
		}
		return fmt.Sprintf("hysteria://%s:%s?%s#%s", n.host, n.port, q.Encode(), name)
	default:
		return ""
	}
}

func ssPayload(cipher, pass string) string {
	if cipher == "" {
		cipher = "aes-256-gcm"
	}
	return base64.StdEncoding.EncodeToString([]byte(cipher + ":" + pass))
}

func base64Encode(raw string) string {
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

func valOr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
