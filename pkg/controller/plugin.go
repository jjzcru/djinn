package controller

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"plugin"

	djinn "github.com/jjzcru/djinn/plugin"
)

func loadPlugins(pluginDirPaths []string) (map[string]djinn.Plugin, error) {
	plugins := make(map[string]djinn.Plugin)
	for _, pluginDirPath := range pluginDirPaths {
		name, plugin, err := loadPlugin(pluginDirPath)
		if err != nil {
			log.Fatal(err)
			continue
			// return plugins, err
		}
		plugins[name] = plugin
	}
	return plugins, nil
}

func loadPlugin(pluginDirPath string) (string, djinn.Plugin, error) {
	modulePath := filepath.Join(pluginDirPath, "plugin.so")
	configPath := filepath.Join(pluginDirPath, "config.yml")
	if _, err := os.Stat(modulePath); os.IsNotExist(err) {
		return "", nil, err
	}

	plug, err := plugin.Open(modulePath)
	if err != nil {
		return "", nil, err
	}

	symPlugin, err := plug.Lookup("Plugin")
	if err != nil {
		return "", nil, err
	}

	p, ok := symPlugin.(djinn.Plugin)
	if !ok {
		return "", nil, errors.New("unexpected type from module symbol")
	}

	err = p.SetConfigFile(configPath)
	if err != nil {
		return "", nil, err
	}

	return filepath.Base(pluginDirPath), p, nil
}

func exists(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		return false
	}

	return true
}
