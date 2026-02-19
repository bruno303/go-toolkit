package config

import (
	"context"
	"embed"
	"os"

	"github.com/bruno303/go-toolkit/pkg/log"
	"github.com/sethvargo/go-envconfig"
	"gopkg.in/yaml.v3"
)

func loadConfig(cfg any, fs embed.FS, loadEnvs bool) {
	filename := "config.yaml"
	if configFile := os.Getenv("CONFIG_FILE"); configFile != "" {
		filename = configFile
	}

	file, err := fs.Open(filename)
	if err != nil {
		panic(err)
	}
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(cfg); err != nil {
		panic(err)
	}

	log.Log().Debug(context.TODO(), "config with yaml: %+v", cfg)

	if !loadEnvs {
		log.Log().Debug(context.TODO(), "skipping env config")
		return
	}

	if err = envconfig.ProcessWith(
		context.Background(),
		&envconfig.Config{
			Target:           cfg,
			DefaultOverwrite: true,
		},
	); err != nil {
		panic(err)
	}

	log.Log().Debug(context.TODO(), "config with envs: %+v", cfg)
}

// LoadConfig and LoadConfigWithoutEnvs load YAML configuration from the provided embed.FS
// and optionally overlay environment variables.
//
// The file to load is taken from the CONFIG_FILE env var or "config.yaml" by default.
// The YAML content is decoded into cfg. When environment processing is enabled (LoadConfig),
// environment variables are applied using github.com/sethvargo/go-envconfig with DefaultOverwrite=true,
// which overwrites struct fields with corresponding env values. Use LoadConfigWithoutEnvs to
// skip environment processing and only load from the YAML file.
//
// Both functions will panic on I/O, decode, or env processing errors.
func LoadConfig(cfg any, fs embed.FS) {
	loadConfig(cfg, fs, true)
}

func LoadConfigWithoutEnvs(cfg any, fs embed.FS) {
	loadConfig(cfg, fs, false)
}
