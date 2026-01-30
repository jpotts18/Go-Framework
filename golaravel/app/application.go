package app

import (
	"fmt"
	"log"
	"net/http"

	"golaravel/app/config"
	"golaravel/app/container"
	"golaravel/app/routing"
	"golaravel/app/view"
)

type Application struct {
	Container *container.Container
	Router    *routing.Router
	Config    *config.Config
	View      *view.View
	basePath  string
}

func New(basePath string) *Application {
	app := &Application{
		Container: container.New(),
		Router:    routing.NewRouter(),
		Config:    config.New(),
		basePath:  basePath,
	}

	app.registerCoreBindings()
	return app
}

func (app *Application) registerCoreBindings() {
	app.Container.Instance("app", app)
	app.Container.Instance("router", app.Router)
	app.Container.Instance("config", app.Config)
}

func (app *Application) SetViewPath(path string) {
	app.View = view.New(path)
	app.Container.Instance("view", app.View)
}

func (app *Application) GetRouter() *routing.Router {
	return app.Router
}

func (app *Application) GetContainer() *container.Container {
	return app.Container
}

func (app *Application) GetConfig() *config.Config {
	return app.Config
}

func (app *Application) GetView() *view.View {
	return app.View
}

func (app *Application) BasePath() string {
	return app.basePath
}

func (app *Application) Run(addr string) error {
	log.Printf("GoLaravel server starting on %s", addr)
	return http.ListenAndServe(addr, app.Router)
}

func (app *Application) RunTLS(addr, certFile, keyFile string) error {
	log.Printf("GoLaravel server starting on %s (TLS)", addr)
	return http.ListenAndServeTLS(addr, certFile, keyFile, app.Router)
}

func (app *Application) Make(abstract string) (interface{}, error) {
	return app.Container.Make(abstract)
}

func (app *Application) MustMake(abstract string) interface{} {
	return app.Container.MustMake(abstract)
}

func (app *Application) Bind(abstract string, concrete container.BindingFunc) {
	app.Container.Bind(abstract, concrete)
}

func (app *Application) Singleton(abstract string, concrete container.BindingFunc) {
	app.Container.Singleton(abstract, concrete)
}

func (app *Application) Instance(abstract string, instance interface{}) {
	app.Container.Instance(abstract, instance)
}

type ServiceProvider interface {
	Register(app *Application)
	Boot(app *Application)
}

func (app *Application) Register(provider ServiceProvider) {
	provider.Register(app)
}

func (app *Application) Boot(provider ServiceProvider) {
	provider.Boot(app)
}

func (app *Application) RegisterProviders(providers []ServiceProvider) {
	for _, provider := range providers {
		app.Register(provider)
	}
}

func (app *Application) BootProviders(providers []ServiceProvider) {
	for _, provider := range providers {
		app.Boot(provider)
	}
}

func (app *Application) Version() string {
	return "GoLaravel 1.0.0"
}

func (app *Application) Environment() string {
	return config.Env("APP_ENV", "development")
}

func (app *Application) IsProduction() bool {
	return app.Environment() == "production"
}

func (app *Application) IsLocal() bool {
	return app.Environment() == "local" || app.Environment() == "development"
}

func (app *Application) Debug() bool {
	return config.EnvBool("APP_DEBUG", true)
}

func (app *Application) URL() string {
	return config.Env("APP_URL", "http://localhost:5000")
}

func PrintBanner() {
	banner := `
   ____       _                                _ 
  / ___| ___ | |    __ _ _ __ __ ___   _____  | |
 | |  _ / _ \| |   / _' | '__/ _' \ \ / / _ \ | |
 | |_| | (_) | |__| (_| | | | (_| |\ V /  __/ | |
  \____|\___/|_____\__,_|_|  \__,_| \_/ \___| |_|
                                                 
  A Laravel-inspired Go Framework
  Version 1.0.0
`
	fmt.Println(banner)
}
