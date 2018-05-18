package hime

import (
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// Config is app's config
type Config struct {
	Globals   map[interface{}]interface{} `yaml:"globals" json:"globals"`
	Routes    map[string]string           `yaml:"routes" json:"routes"`
	Templates []struct {
		Dir        string              `yaml:"dir" json:"dir"`
		Root       string              `yaml:"root" json:"root"`
		Minify     bool                `yaml:"minify" json:"minify"`
		Components []string            `yaml:"components" json:"components"`
		List       map[string][]string `yaml:"list" json:"list"`
		Delims     []string            `yaml:"delims" json:"delims"`
	} `yaml:"templates" json:"templates"`
	Server struct {
		ReadTimeout       string `yaml:"readTimeout" json:"readTimeout"`
		ReadHeaderTimeout string `yaml:"readHeaderTimeout" json:"readHeaderTimeout"`
		WriteTimeout      string `yaml:"writeTimeout" json:"writeTimeout"`
		IdleTimeout       string `yaml:"idleTimeout" json:"idleTimeout"`
		GracefulShutdown  *struct {
			Timeout string `yaml:"timeout" json:"timeout"`
			Wait    string `yaml:"wait" json:"wait"`
		} `yaml:"gracefulShutdown" json:"gracefulShutdown"`
	} `yaml:"server" json:"server"`
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

// Config merges config into app's config
//
// Example:
//
// globals:
//   data1: test
// routes:
//   index: /
//   about: /about
// templates:
// - dir: view
//   root: layout
//   delims: ["{{", "}}"]
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
//   gracefulShutdown:
//     timeout: 1m
//     wait: 5s
func (app *App) Config(config Config) *App {
	app.Globals(config.Globals)
	app.Routes(config.Routes)

	for _, cfg := range config.Templates {
		tp := app.Template()
		tp.Dir(cfg.Dir)
		tp.Root(cfg.Root)
		if len(cfg.Delims) == 2 {
			tp.Delims(cfg.Delims[0], cfg.Delims[1])
		}
		tp.Component(cfg.Components...)
		for name, filenames := range cfg.List {
			tp.Parse(name, filenames...)
		}
		if cfg.Minify {
			tp.Minify()
		}
	}

	// load server config
	parseDuration(config.Server.ReadTimeout, &app.ReadTimeout)
	parseDuration(config.Server.ReadHeaderTimeout, &app.ReadHeaderTimeout)
	parseDuration(config.Server.WriteTimeout, &app.WriteTimeout)
	parseDuration(config.Server.IdleTimeout, &app.IdleTimeout)

	// load graceful config
	if config.Server.GracefulShutdown != nil {
		if app.gracefulShutdown == nil {
			app.gracefulShutdown = &gracefulShutdown{}
		}
		parseDuration(config.Server.GracefulShutdown.Timeout, &app.gracefulShutdown.timeout)
		parseDuration(config.Server.GracefulShutdown.Wait, &app.gracefulShutdown.wait)
	}

	return app
}

// ParseConfig parses config data
func (app *App) ParseConfig(data []byte) *App {
	var config Config
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	return app.Config(config)
}

// ParseConfigFile parses config from file
func (app *App) ParseConfigFile(filename string) *App {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return app.ParseConfig(data)
}
