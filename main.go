package main

import (
	"TZ/controller"
)

func main() {
	ser := controller.NewMyServer()
	ser.Start()
}
