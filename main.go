package main

import (
	"context"
	"log"

	"github.com/mayuka-c/govmomi-practice/config"
	"github.com/mayuka-c/govmomi-practice/internal/client"
)

var serviceConfig config.ServiceConfig

func init() {
	serviceConfig = config.GetServiceConfig(context.Background())
}

func main() {

	client := client.NewVCentreClient(serviceConfig)

	log.Printf("VCenter Version: %s", client.GetVersion())
}
