package main

import (
	"fmt"
	"os"

	"github.com/massalabs/thyra-plugin-wallet/pkg/plugin"
)

//go:generate go run gen.go

func main() {
	// Create a file
	f, err := os.Create("web/html/constants.js")
	if err != nil {
		fmt.Println(err)
		return
	}

	content := "// auto generated file by gen.go\n"
	content += fmt.Sprintf("const pluginAuthor = '%s'; const pluginName = '%s';", plugin.PluginAuthor, plugin.PluginName)

	// Write the content of the file
	_, err = f.WriteString(content)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println("File written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
