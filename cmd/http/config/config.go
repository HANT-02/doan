package config

type RestServer struct {
	Path string `json:"path" yaml:"path"`
	Port string `json:"port" yaml:"port"`
}
