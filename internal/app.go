package tamago

import (
	"fmt"
	"os"
	"path/filepath"
	"tamago/internal/app"

	"github.com/spf13/viper"
)

// Run reads configs and starts the Tamago app
func Run() {
	goPath := os.Getenv("GOPATH")
	envVar := os.Getenv("ENV")

	if envVar == "" {
		envVar = "development"
	}

	viper.Set("ENV", envVar)
	viper.Set("TEMPLATE_PATH", filepath.Join(goPath, "src/force/web/templates/*"))

	viper.SetConfigName(envVar + "_config")
	viper.AddConfigPath(filepath.Join(goPath, "src/force/configs/"))

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	config := viper.GetViper()

	server := app.Server{Config: config}
	server.Init()
	server.Start()
}
