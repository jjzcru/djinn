package controller

import (
	"fmt"
	
	"io/ioutil"
	"os"
	"log"
	"path/filepath"
	"sync"

	"github.com/jjzcru/djinn/plugin"
	nats "github.com/nats-io/nats.go"
)

type Controller struct {
	Plugins map[string]plugin.Plugin
}

func (c *Controller) OnRun() error {
	var wg sync.WaitGroup
	wg.Add(1)
	
	q, err := NewQueue(c.Plugins)
	if err != nil {
		return err
	}
	
	go forever(q)

	wg.Wait()
	q.nc.Close()
	fmt.Println("I finish the OnRun")
	return nil
}

func forever(q queue) {
	fmt.Println("I'm in the forever loop")
	_, err := q.nc.Subscribe("updates", func(m *nats.Msg) {})
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Controller) OnLoadPlugins() error {
	pluginDir := os.Getenv("DJINN_PLUGIN_DIR_PATH")
	files, err := ioutil.ReadDir(pluginDir)
	if err != nil {
		return err
	}

	var pluginDirPaths []string
	for _, f := range files {
		// Validate that the plugin has the required files
		pluginDir := filepath.Join(pluginDir, f.Name())
		pluginFilePath := filepath.Join(pluginDir, "plugin.so")
		infoFilePath := filepath.Join(pluginDir, "info.yml")
		configFilePath := filepath.Join(pluginDir, "config.yml")
		if !exists(pluginFilePath) || !exists(infoFilePath) || !exists(configFilePath) {
			continue
		}
		pluginDirPaths = append(pluginDirPaths, pluginDir)
	}

	plugins, err := loadPlugins(pluginDirPaths)
	if err != nil {
		return err
	}
	c.Plugins = plugins
	return nil
}
