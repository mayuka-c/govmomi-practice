package main

import (
	"context"
	"fmt"
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

	client.GetAllDataCenters()

	virtualMachines, vmTemplates := client.GetVMs()
	log.Printf("Virtual Machine Count: %d", len(virtualMachines))
	fmt.Println("----------------------------")
	for _, vm := range virtualMachines {
		fmt.Println(vm.Name)
	}
	fmt.Printf("\n\n")
	log.Printf("VM Templates Count: %d", len(vmTemplates))
	fmt.Println("----------------------------")
	for _, vmTemplate := range vmTemplates {
		fmt.Println(vmTemplate.Name)
	}
}
