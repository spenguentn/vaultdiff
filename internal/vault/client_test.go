package vault

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: Config{
				Address: "http://localhost:8200",
				Token:   "test-token",
			},
			wantErr: false,
		},
		{
			name: "missing address",
			config: Config{
				Token: "test-token",
			},
			wantErr: true,
		},
		{
			name: "missing token",
			config: Config{
				Address: "http://localhost:8200",
			},
			wantErr: true,
		},
		{
			name:    "empty config",
			config:  Config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client")
			}
			if !tt.wantErr && client.api == nil {
				t.Error("NewClient() returned client with nil api")
			}
		})
	}
}

func TestConfig(t *testing.T) {
	cfg := Config{
		Address: "http://vault.example.com",
		Token:   "s.1234567890abcdef",
	}

	if cfg.Address != "http://vault.example.com" {
		t.Errorf("Config.Address = %v, want %v", cfg.Address, "http://vault.example.com")
	}

	if cfg.Token != "s.1234567890abcdef" {
		t.Errorf("Config.Token = %v, want %v", cfg.Token, "s.1234567890abcdef")
	}
}
