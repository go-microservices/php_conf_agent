package util

import (
	"io/ioutil"
	"log"

	"github.com/GanymedeNil/go-structure-to-php-array"
)

func Write(name string, data map[string]string) {
	arrayData := "<?php\n\nreturn "
	arrayData += go_structure_to_php_array.StructTOPhpArray(data)
	err := ioutil.WriteFile(name+".php", []byte(arrayData), 0777)
	if err != nil {
		log.Fatal(err)
	}
}
