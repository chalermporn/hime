package hime

import (
	"html/template"
	"path/filepath"
)

// TemplateRoot sets root layout using t.Lookup,
// default is "layout"
func (app *App) TemplateRoot(name string) *App {
	app.templateRoot = name
	return app
}

// TemplateDir sets directory to load template,
// default is "view"
func (app *App) TemplateDir(path string) *App {
	app.templateDir = path
	return app
}

// TemplateFuncs adds template funcs while load template
func (app *App) TemplateFuncs(funcs ...template.FuncMap) *App {
	app.templateFuncs = append(app.templateFuncs, funcs...)
	return app
}

// Component adds given templates to every templates
func (app *App) Component(filename ...string) *App {
	app.templateComponents = append(app.templateComponents, filename...)
	return app
}

// Template loads template into memory
func (app *App) Template(name string, filename ...string) *App {
	if _, ok := app.template[name]; ok {
		panic(newErrTemplateDuplicate(name))
	}

	t := template.New("")

	t.Funcs(template.FuncMap{
		"templateName": func() string {
			return name
		},
		"route":  app.Route,
		"global": app.Global,
		"param": func(name string, value interface{}) *Param {
			return &Param{Name: name, Value: value}
		},
	})

	// register funcs
	for _, fn := range app.templateFuncs {
		t.Funcs(fn)
	}

	// load templates and components
	fn := make([]string, len(filename))
	copy(fn, filename)
	fn = append(fn, app.templateComponents...)
	t = template.Must(t.ParseFiles(joinTemplateDir(app.templateDir, fn...)...))
	t = t.Lookup(app.templateRoot)

	app.template[name] = t

	return app
}

func joinTemplateDir(dir string, filenames ...string) []string {
	xs := make([]string, len(filenames))
	for i, filename := range filenames {
		xs[i] = filepath.Join(dir, filename)
	}
	return xs
}
