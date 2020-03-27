package disks

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

func init() {
	ctx := context.Background()

	c, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	computeService, err := compute.New(c)
	if err != nil {
		log.Fatal(err)
	}

	project := os.Getenv("GCP_PROJECT")
	zonesReq := computeService.Zones.List(project)
	if err := zonesReq.Pages(ctx, func(page *compute.ZoneList) error {
		for _, zone := range page.Items {
			disksReq := computeService.Disks.List(project, zone.Name)
			if err := disksReq.Pages(ctx, func(page *compute.DiskList) error {
				for _, disk := range page.Items {
					fmt.Printf("name: %v status: %v zone: %v\n", disk.Name, disk.Status, zone.Name)
				}
				return nil
			}); err != nil {
				log.Fatal(err)
			}
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}
}
