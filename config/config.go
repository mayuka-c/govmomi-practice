package config

import (
	"context"
	"log"

	"github.com/kelseyhightower/envconfig"
)

type ServiceConfig struct {
	VCenterIP string `envconfig:"VCENTER_IP" required:"true"`
	Username  string `envconfig:"VCENTER_USERNAME" required:"true"`
	Password  string `envconfig:"VCENTER_PASSWORD" required:"true"`
}

// GetServiceConfig method to fetch the ServiceConfig
func GetServiceConfig(ctx context.Context) ServiceConfig {
	log.Println("Fetching Service configs")
	config := ServiceConfig{}

	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatalln(ctx, "Failed fetching service configs")
		panic(err)
	}
	return config
}
