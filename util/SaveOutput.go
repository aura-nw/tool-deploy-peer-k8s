package util

import (
	"fmt"
	"log"
	"os"
	"proj1/model"
	"text/template"
	"time"
)

func SaveOutput(t template.Template, fileName string, config model.Config) (string, error) {

	filePath := fmt.Sprintf("./output/%d", time.Now().Unix())
	os.Mkdir(filePath, 0777)
	f, err := os.Create(fmt.Sprintf("%s/%s", filePath, fileName))
	if err != nil {
		log.Println("create file: ", err)
		return "", err
	}
	err = t.Execute(f, config)
	if err != nil {
		log.Print("execute: ", err)
		return "", err
	}
	return filePath, nil
}
