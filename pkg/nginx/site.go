package nginx

import (
	"fmt"
	"os"
	"strings"
)

type Site struct {
	Available     []string
	Enabled       []string
	AvailablePath string
	EnabledPath   string
}

func (s *Site) Enable(name string) {
	var site string
	for _, s := range s.Available {
		if s == strings.TrimSpace(name) {
			site = name
		}
	}
	if site == "" {
		fmt.Fprintln(os.Stderr, "Site not found!")
		os.Exit(1)
	}

	err := os.Symlink(s.AvailablePath+site, s.EnabledPath+site)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not create symlink: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Site %s is enabled!\n", site)
}

func (s *Site) Disable(name string) {
	var site string
	for _, s := range s.Enabled {
		if s == strings.TrimSpace(name) {
			site = name
		}
	}
	if site == "" {
		fmt.Fprintln(os.Stderr, "Site not found!")
		os.Exit(1)
	}
	err := os.Remove(s.EnabledPath + site)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to remove %s: %v\n", site, err)
		os.Exit(1)
	}
	fmt.Println("Site disabled!")
}

func (s *Site) List(list string) {
	switch list {
	case "available":
		fmt.Printf("\nAvailable sites: %v\n\n", s.Available)
	case "enabled":
		fmt.Printf("\nEnabled sites: %v\n\n", s.Enabled)
	case "all":
		fmt.Printf("\nAvailable sites: %v\n", s.Available)
		fmt.Printf("Enabled sites: %v\n\n", s.Enabled)
	default:
		fmt.Fprintln(os.Stderr, "Invalid supplient! Valid values are: available|enabled")
		os.Exit(1)
	}
}
