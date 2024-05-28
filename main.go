package main

import "fmt"

func main() {
	roles, err := getAnsibleRoles("./ansible")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(roles)
}
