package plugin

type Plugin interface {
	SetConfigFile(configFile string) error
	Discover() ([]map[string]interface{}, error)
	Command(id string, command string, payload string) error
}
