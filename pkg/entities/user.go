package entities

type UserPasswd struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Uid      int    `yaml:"uid"`
	Gid      int    `yaml:"gid"`
	Info     string `yaml:"info"`
	Homedir  string `yaml:"homedir"`
	Shell    string `yaml:"shell"`
}
