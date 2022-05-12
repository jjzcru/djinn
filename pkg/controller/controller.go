package controller

import (
	"fmt"
	"log"

	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"encoding/json"

	"github.com/jjzcru/djinn/plugin"
	nats "github.com/nats-io/nats.go"
)

type Controller struct {
	Plugins map[string]plugin.Plugin
	queue queue
}

func (c *Controller) OnRun() error {
	var wg sync.WaitGroup
	wg.Add(1)
	
	q, err := NewQueue(c.Plugins)
	if err != nil {
		return err
	}
	
	c.queue = q
	
	go c.run()

	wg.Wait()
	q.nc.Close()
	fmt.Println("I finish the OnRun")
	return nil
}

func (c *Controller) run() {
	fmt.Println("I'm in the forever loop")
	c.queue.nc.Subscribe("command", c.onReceiveCommand)
}

func (c *Controller) onReceiveCommand(m *nats.Msg) {
	fmt.Printf("Received a message: %s\n", string(m.Data))
	command := Command{}
	if err := json.Unmarshal(m.Data, &command); err != nil {
		fmt.Println("I screw up")
		log.Fatal(err)
		return
	}
	
	fmt.Printf("\nId: %s", command.Id);
	fmt.Printf("\nPlugin: %s", command.Plugin);
	fmt.Printf("\nCommand: %s", command.Command);
	fmt.Printf("\nPayload: %s", command.Payload);
	
	if plugin, ok := c.Plugins[command.Plugin]; ok {
		plugin.Command(command.Id, command.Command, command.Payload)
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
