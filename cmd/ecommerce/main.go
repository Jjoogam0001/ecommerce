package main

import (
	"context"

	"dev.azure.com/jjoogam0290/HelloWorld/HelloWorld/api"
	"dev.azure.com/jjoogam0290/HelloWorld/HelloWorld/config"
	"github.com/labstack/gommon/log"
)

func main() {
	c, err := config.BuildHost().ReadAppConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	pool, err := c.StartDatabase(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	if err != nil {
		log.Error(err)
	}

	if err := api.NewAPI(c).UsingDefaultControllers(pool).Start(c.HTTPServer.Addr); err != nil {
		log.Fatal(err)
	}
}
