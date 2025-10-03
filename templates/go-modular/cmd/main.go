package main

import (
	"fmt"
	"os"

	"{{ package_name | kebab_case }}/cmd/commands"
)

// @title	    Go Application API
// @description	Go Application API documentation
// @version		1.0

// @securityDefinitions.http bearerAuth
// @scheme bearer
// @bearerFormat JWT

// @BasePath /api
func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
