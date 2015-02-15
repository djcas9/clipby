package main

import (
	"io/ioutil"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/Unknwon/com"
)

type Plugin struct {
	Name string
	Path string
	Code string
}

func GetPlugins(pluginPath string) ([]*Plugin, error) {

	var plugins []*Plugin
	files, err := com.GetFileListBySuffix(pluginPath, ".rb")

	if err != nil {
		return plugins, err
	}

	for _, file := range files {

		plugin := &Plugin{
			Name: filepath.Base(file),
			Path: path.Join(PluginPath, file),
		}

		code, err := ioutil.ReadFile(plugin.Path)

		if err != nil {
			log.Warnf("Failed to load plugin: %s", plugin.Name)
		} else {
			plugin.Code = string(code)
			plugins = append(plugins, plugin)
		}
	}

	return plugins, nil
}
