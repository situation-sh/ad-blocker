// request_test.go
package main

import (
	"fmt"
	"os/exec"
	"testing"
)

func TestMain(t *testing.T) {
	fmt.Println(" ####### START #######")
	fmt.Println(" ")
	fmt.Println(" ============ SERVER ============")
	go main()
	// cmd := exec.Command("dig", "@localhost -p 8090 google.com")
	out, _ := exec.Command("dig", "@localhost", "-p8053", "google.com").Output()
	fmt.Println(" ================================")
	fmt.Println(" ")
	fmt.Println(" ")
	fmt.Println(" ")
	fmt.Println(" ============ CLIENT ============")
	fmt.Println(string(out))
	fmt.Println(" ================================")
	fmt.Println(" ")
	fmt.Println(" ")
	fmt.Println(" ")
	fmt.Println(" ####### END #######")
	// send DNS packet to this server
}
