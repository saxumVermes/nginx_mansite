package nginx

import (
	"fmt"
	"os"
	"strings"
)

type Info struct {
	SitesAvailable     []string
	SitesEnabled       []string
	SitesAvailablePath string
	SitesEnabledPath   string
}

func (i *Info) EnableSite(name string) {
	var site string
	for _, s := range i.SitesAvailable {
		if s == strings.TrimSpace(name) {
			site = name
		}
	}
	if site == "" {
		fmt.Fprintln(os.Stderr, "Site not found!")
		os.Exit(1)
	}

	err := os.Symlink(i.SitesAvailablePath+site, i.SitesEnabledPath+site)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not create symlink: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Site %s is enabled!\n", site)
}

func (i *Info) DisableSite(name string) {
	var site string
	for _, s := range i.SitesEnabled {
		if s == strings.TrimSpace(name) {
			site = name
		}
	}
	if site == "" {
		fmt.Fprintln(os.Stderr, "Site not found!")
		os.Exit(1)
	}
	err := os.Remove(i.SitesEnabledPath + site)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to remove %s: %v\n", site, err)
		os.Exit(1)
	}
	fmt.Println("Site disabled!")
}

func (i *Info) ListSites(list string) {
	switch list {
	case "available":
		fmt.Printf("\nAvailable sites: %v\n\n", i.SitesAvailable)
	case "enabled":
		fmt.Printf("\nEnabled sites: %v\n\n", i.SitesEnabled)
	case "all":
		fmt.Printf("\nAvailable sites: %v\n", i.SitesAvailable)
		fmt.Printf("Enabled sites: %v\n\n", i.SitesEnabled)
	default:
		fmt.Fprintln(os.Stderr, "Invalid supplient! Valid values are: available|enabled")
		os.Exit(1)
	}
}
