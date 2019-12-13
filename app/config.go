package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

//Data данные приложения
type Data struct {
	logging *log.Logger
	Cancel  chan os.Signal
}

//GetPath получить путь до запускаемого файла
func GetPath() string {
	if strings.Index(os.Args[0], "debug") > -1 {
		return "debug"
	}
	name := strings.Replace(os.Args[0], ".go", "", -1)
	name = strings.Replace(name, ".exe", "", -1)

	return name
}

//Load загрузка настроек из конфигурационного файла
//Имя конфигурационного файла совпадает с именем запускаемого файла
//Расширение конфигурационого файла .conf
func Load(field ...string) map[interface{}]interface{} {
	fileConfig := fmt.Sprintf("%s.conf", GetPath())

	if _, err := os.Stat(fileConfig); os.IsNotExist(err) {
		log.Println("Config file is not exist")
		return nil
	}

	file, err := ioutil.ReadFile(fileConfig)
	if err != nil {
		log.Fatalln(err)
	}

	var cfg map[interface{}]interface{}
	if err = yaml.Unmarshal(file, &cfg); err != nil {
		log.Println(err)
		return nil
	}

	if len(field) == 0 {
		return cfg
	}

	return cfg[field[0]].(map[interface{}]interface{})
}

//Init установка начальных соединений
func Init(args ...string) (conn Data) {
	conn = Data{
		logging: openLog(),
		Cancel:  make(chan os.Signal),
	}

	signal.Notify(conn.Cancel, os.Kill)
	conn.Log("Start\n")

	return
}

//Start ждет прерывания работы
func (conn Data) Start() {
	<-conn.Cancel
	conn.Stop()
}

//Stop закрывает все соединения
func (conn *Data) Stop() {
	conn.Log("Stop\n")
}
