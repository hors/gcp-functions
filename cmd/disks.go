package disks

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/compute/v1"
)

// CleanDisks : remeve unused disks
func CleanDisks(http.ResponseWriter, *http.Request) {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	project := os.Getenv("GCP_DEV_PROJECT")
	zonesReq := computeService.Zones.List(project)
	if err := zonesReq.Pages(ctx, func(page *compute.ZoneList) error {
		for _, zone := range page.Items {
			disksReq := computeService.Disks.List(project, zone.Name)
			if err := disksReq.Pages(ctx, func(page *compute.DiskList) error {
				for _, disk := range page.Items {
					if disk.Users == nil {
						resp, err := computeService.Disks.Delete(project, zone.Name, disk.Name).Context(ctx).Do()
						if err != nil {
							log.Fatal(err)
						}
						fmt.Printf("%#v\n", resp.Status)
						fmt.Printf("disk: %v in zone %v was deleted\n", disk.Name, zone.Name)
					}
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
