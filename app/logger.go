package app

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

var logName string = fmt.Sprintf("%s.log", GetPath())

func openLog() *log.Logger {
	e, err := os.OpenFile(logName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	logger := log.New(e, "", log.Ldate|log.Ltime)
	logger.SetOutput(&lumberjack.Logger{
		Filename:   logName,
		MaxSize:    30,
		MaxBackups: 0,
		MaxAge:     30,
	})

	return logger
}

//Log пишет в файл лог файл
func (app *Data) Log(mess string, arg ...interface{}) {
	mess = fmt.Sprintf(mess, arg...)

	fmt.Println(mess)

	app.logging.SetPrefix("")
	app.logging.Println(mess)

}

//Err логирование ошибок
func (app *Data) Err(mess string, arg ...interface{}) {
	mess = fmt.Sprintf(mess, arg...)

	fmt.Println(mess)

	app.logging.SetPrefix("ERROR ")
	app.logging.Println(mess)
}

//Completion write info about ending function
func (app Data) Completion(mess string, arg ...interface{}) {
	if r := recover(); r != nil {
		app.Err(mess, r.(string))
	} else {
		app.Log("End %s", fmt.Sprintf(mess, arg...))
	}
}
