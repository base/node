package config

import "testing"

func TestLoadDefaultConfig(t *testing.T) {
    cfg := DefaultConfig()
    if cfg.HTTPPort == 0 {
        t.Fatalf("HTTPPort should be set in default config")
    }
    if cfg.DataDir == "" {
        t.Fatalf("DataDir should not be empty in default config")
    }
}
