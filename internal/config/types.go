package config

type Instruction struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Spec       any    `yaml:"spec"`
}
