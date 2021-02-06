package main

import (
	"../../ccDoc"
	"os"
)

func main() {
	ccDoc.GenerateDocJson( os.Args[1], os.Args[2] )
}
