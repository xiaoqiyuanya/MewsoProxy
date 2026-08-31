package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
)

const (
	secretsAccessKey  = "MEWSO_JWT_ACCESS_SECRET"
	secretsRefreshKey = "MEWSO_JWT_REFRESH_SECRET"
	secretsSSHKey     = "MEWSO_APP_SSH_ENCRYPT_KEY"
	secretsTokenKey   = "MEWSO_APP_SERVER_TOKEN"
)

func EnsureSecrets(cfg *Config, path string) error {
	if path == "" {
		generateIfNeeded(cfg)
		return nil
	}
	if vals, err := readSecretsFile(path); err == nil {
		cfg.JWT.AccessSecret = vals[secretsAccessKey]
		cfg.JWT.RefreshSecret = vals[secretsRefreshKey]
		cfg.App.SSHEncryptKey = vals[secretsSSHKey]
		cfg.App.ServerToken = vals[secretsTokenKey]
	}
	changed := generateIfNeeded(cfg)
	if changed || !fileExists(path) {
		if err := writeSecretsFile(path, cfg); err != nil {
			return fmt.Errorf("write secrets file: %w", err)
		}
	}
	return nil
}

func generateIfNeeded(cfg *Config) bool {
	changed := false
	if needsGen(cfg.JWT.AccessSecret) {
		cfg.JWT.AccessSecret = randHex(32)
		changed = true
	}
	if needsGen(cfg.JWT.RefreshSecret) {
		cfg.JWT.RefreshSecret = randHex(32)
		changed = true
	}
	if needsGen(cfg.App.SSHEncryptKey) {
		cfg.App.SSHEncryptKey = randHex(32)
		changed = true
	}
	if needsGen(cfg.App.ServerToken) {
		cfg.App.ServerToken = randHex(32)
		changed = true
	}
	return changed
}

func needsGen(v string) bool {
	if v == "" {
		return true
	}
	lv := strings.ToLower(v)
	return strings.Contains(lv, "change-me") || strings.Contains(lv, "please-change")
}

func randHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("crypto/rand failed: %v", err))
	}
	return hex.EncodeToString(b)
}

func readSecretsFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	out := map[string]string{}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		out[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	return out, nil
}

func writeSecretsFile(path string, cfg *Config) error {
	content := "# 首次启动自动生成的密钥，请勿泄漏；如需固定可手工编辑本文件\n" +
		secretsAccessKey + "=" + cfg.JWT.AccessSecret + "\n" +
		secretsRefreshKey + "=" + cfg.JWT.RefreshSecret + "\n" +
		secretsSSHKey + "=" + cfg.App.SSHEncryptKey + "\n" +
		secretsTokenKey + "=" + cfg.App.ServerToken + "\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		return err
	}
	return os.Chmod(path, 0600)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
