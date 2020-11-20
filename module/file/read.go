package file

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

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

func FolderFiles() []string {
	//取得執行檔的路徑,包含檔名
	file, _ := exec.LookPath(os.Args[0])

	//得到全路径，比如在windows下E:\\golang\\test\\a.exe
	filePath, _ := filepath.Abs(file)
	//只取資料夾
	pwd := filepath.Dir(filePath)
	fmt.Println("pwd", pwd)
	files, err := ioutil.ReadDir(pwd)
	if err != nil {
		log.Fatal(err)
	}
	filesName := []string{}
	for _, file := range files {
		filesName = append(filesName, file.Name())
	}

	return filesName
}
