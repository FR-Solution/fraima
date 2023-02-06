package config

type Instruction struct {
	Metadata `yaml:",inline"`
	Spec     Spec `yaml:"spec"`
}

type Metadata struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
}

type Spec struct {
	Service       *Config               `yaml:"service,omitempty"`
	Configuration *Config               `yaml:"configuration,omitempty"`
	Download      []DownloadInstruction `yaml:"download"`
	Starting      []string              `yaml:"starting"`
}

type Config struct {
	ExtraArgs any `yaml:"extraArgs"`
}

type DownloadInstruction struct {
	Name       string    `yaml:"name"`
	Src        string    `yaml:"src"`
	CheckSum   *CheckSum `yaml:"checkSum"`
	HostPath   string    `yaml:"path"`
	Owner      string    `yaml:"owner"`
	Permission int       `yaml:"permission"`
	Unzip      Unzip     `yaml:"unzip"`
}

type Unzip struct {
	Status bool     `yaml:"status"`
	Files  []string `yaml:"files"`
}

type CheckSum struct {
	Src  string `yaml:"src"`
	Type string `yaml:"type"`
}
