package main

import (
	"context"
	"dev.azure.com/jjoogam/Ecommerce-core/config"
	"log"

	"dev.azure.com/jjoogam/Ecommerce-core/api"
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
		log.Fatal(err)
	}

	if err := api.NewAPI(c).UsingDefaultControllers(pool).Start(c.HTTPServer.Addr); err != nil {
		log.Fatal(err)
	}
}
