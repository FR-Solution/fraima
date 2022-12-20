package config

type Config struct {
	Files []File `yaml:",inline"`
}

type File struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	ExtraArgs  any    `yaml:"extraArgs"`
}
