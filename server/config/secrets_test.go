package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestEnsureSecretsGeneratesInMemory(t *testing.T) {
	cfg := &Config{}
	cfg.JWT.AccessSecret = "please-change-access-secret"
	cfg.JWT.RefreshSecret = "change-me-refresh"
	cfg.App.SSHEncryptKey = "change-me-ssh-encrypt-key-32bytes"
	cfg.App.ServerToken = "change-me-server-token"
	if err := EnsureSecrets(cfg, ""); err != nil {
		t.Fatalf("EnsureSecrets error: %v", err)
	}
	if cfg.JWT.AccessSecret == "" || strings.Contains(cfg.JWT.AccessSecret, "change-me") || strings.Contains(cfg.JWT.AccessSecret, "please-change") {
		t.Fatalf("access secret not regenerated: %q", cfg.JWT.AccessSecret)
	}
	if cfg.JWT.RefreshSecret == "" || cfg.App.SSHEncryptKey == "" || cfg.App.ServerToken == "" {
		t.Fatal("secrets should be generated")
	}
}

func TestEnsureSecretsPersistAndReload(t *testing.T) {
	path := filepath.Join(t.TempDir(), "secrets.env")
	cfg := &Config{}
	cfg.JWT.AccessSecret = "change-me-access"
	cfg.JWT.RefreshSecret = "change-me-refresh"
	cfg.App.SSHEncryptKey = "change-me-ssh"
	cfg.App.ServerToken = "change-me-token"
	if err := EnsureSecrets(cfg, path); err != nil {
		t.Fatalf("EnsureSecrets write error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("secrets file not written: %v", err)
	}
	if runtime.GOOS != "windows" {
		info, _ := os.Stat(path)
		if info.Mode().Perm()&077 != 0 {
			t.Fatalf("secrets file perms should be 0600, got %v", info.Mode().Perm())
		}
	}

	cfg2 := &Config{}
	cfg2.JWT.AccessSecret = "change-me-access"
	cfg2.JWT.RefreshSecret = "change-me-refresh"
	cfg2.App.SSHEncryptKey = "change-me-ssh"
	cfg2.App.ServerToken = "change-me-token"
	if err := EnsureSecrets(cfg2, path); err != nil {
		t.Fatalf("EnsureSecrets reload error: %v", err)
	}
	if cfg2.JWT.AccessSecret != cfg.JWT.AccessSecret || cfg2.JWT.RefreshSecret != cfg.JWT.RefreshSecret ||
		cfg2.App.SSHEncryptKey != cfg.App.SSHEncryptKey || cfg2.App.ServerToken != cfg.App.ServerToken {
		t.Fatal("reload should reuse persisted secrets")
	}
}

func TestEnsureSecretsKeepsExplicitValues(t *testing.T) {
	cfg := &Config{}
	cfg.JWT.AccessSecret = "my-custom-secret-abcdef"
	cfg.JWT.RefreshSecret = "my-custom-refresh-abcdef"
	cfg.App.SSHEncryptKey = "my-custom-ssh-abcdef"
	cfg.App.ServerToken = "my-custom-token-abcdef"
	if err := EnsureSecrets(cfg, ""); err != nil {
		t.Fatalf("EnsureSecrets error: %v", err)
	}
	if cfg.JWT.AccessSecret != "my-custom-secret-abcdef" {
		t.Fatalf("explicit secret should be kept, got %q", cfg.JWT.AccessSecret)
	}
}
