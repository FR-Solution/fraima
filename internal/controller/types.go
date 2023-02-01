package controller

type GenerateInstruction any

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
