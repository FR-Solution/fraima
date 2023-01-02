package config

type Config struct {
	GenerateList []Generate `yaml:"generating"`
	DownloadList []Download `yaml:"downloading"`
}

type Generate struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	ExtraArgs  any    `yaml:"extraArgs"`
}

type Download struct {
	URL        string `yaml:"url"`
	Filepath   string `yaml:"filepath"`
	Permission int    `yaml:"permission"`
}
