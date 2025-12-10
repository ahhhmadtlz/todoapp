package config

import (
	"log"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func Load(configPath string) Config {
	var k = koanf.New(".")

	// Load default config
	if err := k.Load(confmap.Provider(defaultConfig, "."), nil); err != nil {
		log.Printf("⚠️  Error loading default config: %v", err)
	}

	// Load file config (should override defaults)
	if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		log.Fatalf("❌ Error loading config file '%s': %v", configPath, err)
	}

	// Load environment variables (highest priority)
	if err := k.Load(env.Provider("TODOAPP_", ".", func(s string) string {
		str := strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "TODOAPP_")), "_", ".", -1)
		return strings.Replace(str, "..", "_", -1)
	}), nil); err != nil {
		log.Printf("⚠️  Error loading env vars: %v", err)
	}

	// Debug: Print what koanf loaded
	log.Printf("🔍 Loaded config keys: %v", k.Keys())
	log.Printf("🔍 MySQL config: username=%s, host=%s, port=%d, db=%s",
		k.String("mysql.username"),
		k.String("mysql.host"),
		k.Int("mysql.port"),
		k.String("mysql.db_name"))

	var cfg Config

	if err := k.Unmarshal("", &cfg); err != nil {
		log.Fatalf("❌ Error unmarshaling config: %v", err)
	}

	return cfg
}