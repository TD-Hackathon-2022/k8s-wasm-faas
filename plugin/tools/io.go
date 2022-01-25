package tools

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func CloseResource(resource io.Closer) {
	err := resource.Close()
	if err != nil {
		fmt.Println("resource close file")
		os.Exit(1)
	}
}

func GetFileName(filepath string) string {
	stat, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		fmt.Printf("faas script file is not existed: %s\n", filepath)
		os.Exit(1)
	}

	return stat.Name()
}

func ReadAll(filepath string) string {

	file, err := os.Open(filepath)
	defer CloseResource(file)

	if err != nil {
		fmt.Println("open faas file error")
		os.Exit(1)
	}

	content, _ := ioutil.ReadAll(file)
	return string(content)
}
