package main

import (
	"flag"

	"github.com/tgascoigne/pogo"
	_ "github.com/tgascoigne/pogo/json"
)

var path = flag.String("path", "", "")
var pkg = flag.String("pkg", "", "")

func main() {
	flag.Parse()
	pogo.VisitorTemplate.Generate("json", *pkg, *path)
}
