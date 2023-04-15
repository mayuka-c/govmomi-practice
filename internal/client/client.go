package client

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/mayuka-c/govmomi-practice/config"
	"github.com/vmware/govmomi"
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
