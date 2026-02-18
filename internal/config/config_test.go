package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Theme != DefaultTheme {
		t.Errorf("Theme = %q, want %q", cfg.Theme, DefaultTheme)
	}
	if cfg.TimestampFormat != DefaultTimestampFormat {
		t.Errorf("TimestampFormat = %q, want %q", cfg.TimestampFormat, DefaultTimestampFormat)
	}
	if cfg.ShowLineNumbers != DefaultShowLineNumbers {
		t.Errorf("ShowLineNumbers = %v, want %v", cfg.ShowLineNumbers, DefaultShowLineNumbers)
	}
	if cfg.WrapLines != DefaultWrapLines {
		t.Errorf("WrapLines = %v, want %v", cfg.WrapLines, DefaultWrapLines)
	}
	if cfg.JSONIndent != DefaultJSONIndent {
		t.Errorf("JSONIndent = %d, want %d", cfg.JSONIndent, DefaultJSONIndent)
	}
	if cfg.MaxBufferSize != DefaultMaxBufferSize {
		t.Errorf("MaxBufferSize = %d, want %d", cfg.MaxBufferSize, DefaultMaxBufferSize)
	}
	if cfg.WorkerCount != DefaultWorkerCount {
		t.Errorf("WorkerCount = %d, want %d", cfg.WorkerCount, DefaultWorkerCount)
	}
	if cfg.Follow {
		t.Error("Follow = true, want false")
	}
}
