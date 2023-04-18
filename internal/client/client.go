package client

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"sort"
	"strings"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/mayuka-c/govmomi-practice/config"
)

type VCentreClient struct {
	client *govmomi.Client
}

func NewVCentreClient(config config.ServiceConfig) *VCentreClient {

	_url, err := url.Parse(fmt.Sprintf("https://%s/sdk", config.VCenterIP))
	if err != nil {
		log.Fatalf("Logging in error: %s\n", err.Error())
		return nil
	}

	_url.User = url.UserPassword(config.Username, config.Password)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := govmomi.NewClient(ctx, _url, true)
	if err != nil {
		log.Fatalf("Logging in error: %s\n", err.Error())
		return nil
	}

	log.Println("Log in to VCenter successful")

	return &VCentreClient{
		client: client,
	}
}

func (vc *VCentreClient) LogoutClient() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	vc.client.Logout(ctx)
}

func (vc *VCentreClient) GetVersion() string {
	return vc.client.ServiceContent.About.Version
}

func (vc *VCentreClient) GetAllDataCenters() ([]*object.Datacenter, error) {
	f := find.NewFinder(vc.client.Client, true)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	datacenters, err := f.DatacenterList(ctx, "*")
	if err != nil {
		log.Fatalf("Logging in error: %s\n", err.Error())
		return nil, err
	}

	return datacenters, nil
}

type ByName []mo.VirtualMachine

func (n ByName) Len() int           { return len(n) }
func (n ByName) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
func (n ByName) Less(i, j int) bool { return strings.ToLower(n[i].Name) < strings.ToLower(n[j].Name) }

func (vc *VCentreClient) GetAllVMs() ([]mo.VirtualMachine, []mo.VirtualMachine) {
	f := find.NewFinder(vc.client.Client, true)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dcList, err := vc.GetAllDataCenters()
	if err != nil {
		log.Fatalf("Logging in error: %s\n", err.Error())
	}

	var virtualMachines, vmTemplates []mo.VirtualMachine

	for _, dc := range dcList {
		// set the current datacenter
		f.SetDatacenter(dc)
		// Find all virtual machines in datacenter
		vms, err := f.VirtualMachineList(ctx, "*")
		if err != nil {
			log.Fatalf("Logging in error: %s\n", err.Error())
		}

		pc := property.DefaultCollector(vc.client.Client)

		// Convert datastores into list of references
		var refs []types.ManagedObjectReference
		for _, vm := range vms {
			refs = append(refs, vm.Reference())
		}
		// Retrieve name property for all vms
		var vmt []mo.VirtualMachine
		err = pc.Retrieve(ctx, refs, nil, &vmt)
		if err != nil {
			log.Fatalf("Logging in error: %s\n", err.Error())
		}

		fmt.Println("Virtual machines found:", len(vmt))
		sort.Sort(ByName(vmt))
		for _, vm := range vmt {
			if vm.Config == nil {
				log.Println("Config is nil")
				continue
			}
			if vm.Config.Template {
				vmTemplates = append(vmTemplates, vm)
			} else {
				virtualMachines = append(virtualMachines, vm)
			}
		}
	}

	return virtualMachines, vmTemplates
}

func (vc *VCentreClient) GetAllDatastores() {
	f := find.NewFinder(vc.client.Client, true)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dcList, err := vc.GetAllDataCenters()
	if err != nil {
		log.Fatalf("Logging in error: %s\n", err.Error())
	}

	var datastoresList []*object.Datastore

	for _, dc := range dcList {
		f.SetDatacenter(dc)
		datastores, err := f.DatastoreList(ctx, "*")
		if err != nil {
			log.Fatalf("Logging in error: %s\n", err.Error())
		}
		datastoresList = append(datastoresList, datastores...)
	}

	fmt.Println("Datastores: ", datastoresList)
}

func (vc *VCentreClient) GetAllResourcePools() {
	f := find.NewFinder(vc.client.Client, true)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dcList, err := vc.GetAllDataCenters()
	if err != nil {
		log.Fatalf("Logging in error: %s\n", err.Error())
	}

	var resourcePoolList []*object.ResourcePool

	for _, dc := range dcList {
		f.SetDatacenter(dc)
		resourcePools, err := f.ResourcePoolList(ctx, "*")
		if err != nil {
			log.Fatalf("Logging in error: %s\n", err.Error())
		}
		resourcePoolList = append(resourcePoolList, resourcePools...)
	}

	fmt.Println("ResourcePools: ", resourcePoolList)
}

func (vc *VCentreClient) GetSnapshotDetailsOfVM(vm mo.VirtualMachine) {
	snapshotList := vm.Snapshot.RootSnapshotList
	fmt.Println("Snapshot: ", snapshotList)
}
