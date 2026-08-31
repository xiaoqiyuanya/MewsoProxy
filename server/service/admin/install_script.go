package admin

import (
	"encoding/json"
	"strconv"
	"strings"
)

const singboxVersion = "v1.11.8"

func buildInstallScript(n installNode, uuid string) string {
	cfg := singboxConfig(n, uuid)
	return strings.Join([]string{
		"#!/usr/bin/env bash",
		"set -e",
		"ARCH=$(uname -m)",
		`case "$ARCH" in x86_64|amd64) ARCH=amd64;; aarch64|arm64) ARCH=arm64;; *) echo "不支持的架构: $ARCH"; exit 1;; esac`,
		"OS=linux",
		`VER="` + singboxVersion + `"`,
		`URL="https://github.com/SagerNet/sing-box/releases/download/$VER/sing-box-$VER-$OS-$ARCH.tar.gz"`,
		`echo "下载 sing-box $VER ($OS/$ARCH)..."`,
		`(command -v wget >/dev/null && wget -qO /tmp/sing.tar.gz "$URL") || curl -fsSL -o /tmp/sing.tar.gz "$URL"`,
		`mkdir -p /opt/sing-box /opt/sing-box/cert`,
		`tar -xzf /tmp/sing.tar.gz -C /tmp`,
		`cp /tmp/sing-box-$VER-$OS-$ARCH/sing-box /opt/sing-box/sing-box`,
		`chmod +x /opt/sing-box/sing-box`,
		`openssl req -x509 -nodes -newkey rsa:2048 -keyout /opt/sing-box/cert/key.pem -out /opt/sing-box/cert/cert.pem -days 3650 -subj "/CN=localhost" 2>/dev/null || true`,
		`cat > /opt/sing-box/config.json <<'CFGEOF'`,
		cfg,
		"CFGEOF",
		`cat > /etc/systemd/system/sing-box.service <<'SVDEOF'`,
		"[Unit]",
		"Description=sing-box",
		"After=network.target",
		"[Service]",
		"ExecStart=/opt/sing-box/sing-box run -c /opt/sing-box/config.json",
		"Restart=on-failure",
		"[Install]",
		"WantedBy=multi-user.target",
		"SVDEOF",
		"systemctl daemon-reload",
		"systemctl enable sing-box",
		"systemctl restart sing-box",
		`echo "sing-box 启动成功，监听端口: ` + portStr(n.serverPort) + `"`,
	}, "\n")
}

func singboxConfig(n installNode, uuid string) string {
	inbound := singboxInbound(n, uuid)
	root := map[string]interface{}{
		"log": map[string]interface{}{"level": "info"},
		"inbounds": []interface{}{inbound},
		"outbounds": []interface{}{
			map[string]interface{}{"type": "freedom"},
		},
	}
	b, _ := json.MarshalIndent(root, "", "  ")
	return string(b)
}

func singboxInbound(n installNode, uuid string) map[string]interface{} {
	base := map[string]interface{}{
		"listen":      "0.0.0.0",
		"listen_port": n.serverPort,
	}
	tls := map[string]interface{}{
		"enabled":         true,
		"certificate_path": "/opt/sing-box/cert/cert.pem",
		"key_path":        "/opt/sing-box/cert/key.pem",
	}
	switch n.nodeType {
	case "trojan":
		base["type"] = "trojan"
		base["users"] = []interface{}{map[string]interface{}{"password": uuid}}
		base["tls"] = tls
	case "vmess":
		base["type"] = "vmess"
		base["users"] = []interface{}{map[string]interface{}{"uuid": uuid, "alterId": 0}}
	case "shadowsocks", "ss":
		base["type"] = "shadowsocks"
		base["method"] = cipherOr(n.cipher)
		base["password"] = uuid
	case "hysteria":
		base["type"] = "hysteria"
		base["auth"] = uuid
		base["up_mbps"] = n.upMbps
		base["down_mbps"] = n.downMbps
		base["tls"] = tls
	}
	return base
}

func cipherOr(c string) string {
	if c == "" {
		return "aes-256-gcm"
	}
	return c
}

func portStr(p int) string {
	if p <= 0 {
		return "0"
	}
	return strconv.Itoa(p)
}
