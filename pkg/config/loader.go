package config

import (
	"context"
	"embed"
	"os"

	"github.com/bruno303/go-toolkit/pkg/log"
	"github.com/sethvargo/go-envconfig"
	"gopkg.in/yaml.v3"
)

// LoadConfig loads the configuration from the given FS and env vars.
//
// It first tries to load the configuration from the file specified by the
// CONFIG_FILE env var, or "config.yaml" if not specified. Then it decodes the
// YAML file into the given cfg struct. After that, it tries to overwrite the
// struct fields with the corresponding env vars.
func LoadConfig(cfg any, fs embed.FS) {
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
