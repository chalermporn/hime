package hime

import (
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// Config is app's config
type Config struct {
	Globals  map[interface{}]interface{} `yaml:"globals" json:"globals"`
	Routes   map[string]string           `yaml:"routes" json:"routes"`
	Template struct {
		Dir        string              `yaml:"dir" json:"dir"`
		Root       string              `yaml:"root" json:"root"`
		Minify     bool                `yaml:"minify" json:"minify"`
		Components []string            `yaml:"components" json:"components"`
		List       map[string][]string `yaml:"list" json:"list"`
	} `yaml:"template" json:"template"`
	Server struct {
		ReadTimeout       string `yaml:"readTimeout" json:"readTimeout"`
		ReadHeaderTimeout string `yaml:"readHeaderTimeout" json:"readHeaderTimeout"`
		WriteTimeout      string `yaml:"writeTimeout" json:"writeTimeout"`
		IdleTimeout       string `yaml:"idleTimeout" json:"idleTimeout"`
	} `yaml:"server" json:"server"`
	Graceful struct {
		Timeout string `yaml:"timeout" json:"timeout"`
		Wait    string `yaml:"wait" json:"wait"`
	} `yaml:"graceful" json:"graceful"`
}

func parseDuration(s string, t *time.Duration) {
	if s == "" {
		return
	}
	var err error
	*t, err = time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
}

// Load loads config
//
// Example:
//
// globals:
//   data1: test
// routes:
//   index: /
//   about: /about
// template:
//   dir: view
//   root: layout
//   minify: true
//   components:
//   - comp/comp1.tmpl
//   - comp/comp2.tmpl
//   list:
//     main.tmpl:
//     - main.tmpl
//     - _layout.tmpl
//     about.tmpl: [about.tmpl, _layout.tmpl]
// server:
//   readTimeout: 10s
//   readHeaderTimeout: 5s
//   writeTimeout: 5s
//   idleTimeout: 30s
// graceful:
//   timeout: 1m
//   wait: 5s
func (app *App) Load(config Config) *App {
	app.Globals(config.Globals)
	app.Routes(config.Routes)
	app.templateDir = config.Template.Dir
	app.templateRoot = config.Template.Root
	app.Component(config.Template.Components...)

	for name, filenames := range config.Template.List {
		app.Template(name, filenames...)
	}

	if config.Template.Minify {
		app.Minify()
	}

	// load server config
	parseDuration(config.Server.ReadTimeout, &app.ReadTimeout)
	parseDuration(config.Server.ReadHeaderTimeout, &app.ReadHeaderTimeout)
	parseDuration(config.Server.WriteTimeout, &app.WriteTimeout)
	parseDuration(config.Server.IdleTimeout, &app.IdleTimeout)

	// load graceful config
	parseDuration(config.Graceful.Timeout, &app.graceful.timeout)
	parseDuration(config.Graceful.Wait, &app.graceful.wait)
	return app
}

// LoadFromFile loads config from file
func (app *App) LoadFromFile(filename string) *App {
	fs, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fs.Close()

	var config Config
	err = yaml.NewDecoder(fs).Decode(&config)
	if err != nil {
		panic(err)
	}

	return app.Load(config)
}