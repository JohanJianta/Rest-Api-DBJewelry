package main

import "Rest-Api-DBJewelry/apiMySql"

func main() {
	// Rest API dengan database file json
	// apiJson.Init()

	// Rest API dengan database mysql di localhost
	apiMySql.Init()
}
