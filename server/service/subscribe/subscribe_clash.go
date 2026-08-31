package subscribe

import (
	"encoding/json"
	"fmt"
	"strings"
)

func vmessJSON(obj map[string]interface{}) string {
	b, _ := json.Marshal(obj)
	return string(b)
}

func buildClash(nodes []nodeInfo, uuid string) string {
	var b strings.Builder
	b.WriteString("proxies:\n")
	for _, n := range nodes {
		b.WriteString(clashProxy(n, uuid))
	}
	b.WriteString("\nproxy-groups:\n")
	b.WriteString("  - name: 🚀 节点选择\n    type: select\n    proxies:\n")
	for _, n := range nodes {
		b.WriteString("      - " + clashProtocol(n.typ) + ": " + n.name + "\n")
	}
	return b.String()
}

func clashProxy(n nodeInfo, uuid string) string {
	switch n.typ {
	case "trojan":
		return fmt.Sprintf("  - name: %s\n    type: trojan\n    server: %s\n    port: %s\n    password: %s\n    sni: %s\n    udp: true\n", n.name, n.host, n.port, uuid, n.sni)
	case "vmess":
		return fmt.Sprintf("  - name: %s\n    type: vmess\n    server: %s\n    port: %s\n    uuid: %s\n    alterId: 0\n    cipher: auto\n    network: %s\n    tls: %t\n    udp: true\n", n.name, n.host, n.port, uuid, n.network, n.tls)
	case "shadowsocks":
		return fmt.Sprintf("  - name: %s\n    type: ss\n    server: %s\n    port: %s\n    cipher: %s\n    password: %s\n    udp: true\n", n.name, n.host, n.port, n.cipher, uuid)
	case "hysteria":
		return fmt.Sprintf("  - name: %s\n    type: hysteria\n    server: %s\n    port: %s\n    auth: %s\n    up: %s\n    down: %s\n    sni: %s\n    skip-cert-verify: %t\n", n.name, n.host, n.port, uuid, mbps(n.upMbps), mbps(n.downMbps), n.sni, n.insecure)
	default:
		return ""
	}
}

func clashProtocol(typ string) string {
	switch typ {
	case "shadowsocks":
		return "SS"
	case "hysteria":
		return "Hysteria"
	default:
		return strings.ToUpper(typ)
	}
}

func mbps(v int) string {
	if v <= 0 {
		return "100"
	}
	return fmt.Sprintf("%d", v)
}
