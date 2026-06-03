package main

import (
	"aggregator/internal/config"
	"fmt"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = cfg.SetUser("superz97")
	if err != nil {
		fmt.Println(err)
		return
	}

	updatedCfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", updatedCfg)
}
