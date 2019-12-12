package file

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type (
	config struct {
		Name string
		Age  int
	}
)

func ReadFile() {
	var c config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("has error", err.Error())
		return
	}
	if err := viper.Unmarshal(&c); err != nil {
		fmt.Println("has error", err.Error())
		return
	}
	fmt.Println("done", c)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("is changed", e.Name)
	})

	fmt.Println("port is", viper.GetInt("port"))
	return

}
