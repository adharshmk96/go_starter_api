/*
Copyright © 2025 Adharsh Manikandan <debugslayer@gmail.com>
*/
package main

import (
	"servicehub_api/cmd"
	_ "servicehub_api/docs"
)

//	@title			ServiceHub API
//	@version		1.0
//	@description	This is the API for the ServiceHub platform.
//	@host			localhost:8080
//	@BasePath		/api/v1
//	@schemes		http
func main() {
	cmd.Execute()
}
