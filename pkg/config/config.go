package config

type StorageConfig struct {
	Root      string `json:"root" yaml:"root"`
	ImagePath string `json:"imagePath" yaml:"imagePath"`
	Domain    string `json:"domain" yaml:"domain"`
	Rec       string `json:"rec" yaml:"rec"`
	Snap      string `json:"snap" yaml:"snap"`
	Ext       string `json:"ext" yaml:"ext"`
}
