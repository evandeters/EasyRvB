package main

import "fmt"

func main() {
	roles, err := getAnsibleRoles("./ansible")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(roles)

	configMap := make(map[string]map[string]string)

	for _, role := range roles {
		configMap, err = parseRoleMetadata(role)
		for k := range configMap {
			fmt.Printf("Category: %s\n", k)
			for k2, v2 := range configMap[k] {
				fmt.Printf("Key: %s ", k2)
				fmt.Printf("Value: %s\n", v2)
			}
		}
	}
}
