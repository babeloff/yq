//+build ignore

// Zero install
// https://magefile.org/zeroinstall/
package main

import (
	"github.com/magefile/mage/mage"
	"os"
)

// go run mage.go
// go run mage.go build
func main() {
	os.Exit(mage.Main())
}
