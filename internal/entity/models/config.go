package models

type Config struct {
	Database struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Name string `yaml:"name"`
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
	} `yaml:"database"`

	Port string `yaml:"port"`

	Org struct {
		OrgType    string `yaml:"orgType"`
		OrgId      string `yaml:"orgId"`
		SignMethod string `yaml:"signMethod"`
		PrivateKey string `yaml:"privateKey"`
	} `yaml:"org"`
}
