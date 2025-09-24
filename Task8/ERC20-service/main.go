package main

import (
	"erc20-service/cmd"
	_ "erc20-service/cmd/backfill"
	_ "erc20-service/cmd/daemon"
	_ "erc20-service/cmd/health"
)

func main() {
	cmd.Execute()
}
