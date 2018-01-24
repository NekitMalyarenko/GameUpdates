package main

import (
	"os"
	"log"
)


func main() {

	for _, val := range os.Environ() {
		log.Println(val)
	}
}