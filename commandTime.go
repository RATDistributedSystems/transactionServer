package main

import(
	"fmt"
	"os"

)


func appendToText(filename string, text string){

	//openfile
	text = text + "\n"
	fmt.Println("Appending To File: " + filename)
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
  	  panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
   	 panic(err)
	}

}

func createFile(filename string) {

	// detect if file exists
	var _, err = os.Stat(filename)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(filename)

		if err!= nil { 
			return 
		}

		defer file.Close()
	}

	fmt.Println("==> done creating file", filename)
}