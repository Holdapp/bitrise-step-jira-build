package config

type JIRAConfig struct {
	Host               string
	Username           string
	Token              string
	ProjectKey         string
	BuildCustomFieldID string
}

func (config *JIRAConfig) Validate() bool {
	return true
}
