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

	virtualMachines, vmTemplates := client.GetAllVMs()
	log.Printf("Virtual Machine Count: %d", len(virtualMachines))
	fmt.Println("----------------------------")
	for _, vm := range virtualMachines {
		fmt.Println(vm.Name)
		if vm.Snapshot != nil {
			fmt.Println("Snapshot")
			fmt.Println("--------------------------")
			client.GetSnapshotDetailsOfVM(vm)
		}
	}
	fmt.Printf("\n\n")
	log.Printf("VM Templates Count: %d", len(vmTemplates))
	fmt.Println("----------------------------")
	for _, vmTemplate := range vmTemplates {
		fmt.Println(vmTemplate.Name)
	}

	fmt.Printf("\n\n")
	client.GetAllDatastores()

	fmt.Printf("\n\n")
	client.GetAllResourcePools()
}
