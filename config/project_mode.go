package config

type ProjectMode struct {
	Env string `envconfig:"PROJECT_MODE" default:"development"`
}
