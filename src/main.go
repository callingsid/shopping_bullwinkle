package main
 import (
 	"fmt"
	 "github.com/callingsid/shopping_bullwinkle/src/app"
 )
var appName = "bullwinkle"

func main() {
	fmt.Printf("Starting %v\n", appName)
	app.StartApplication()
}

