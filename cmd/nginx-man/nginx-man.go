package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/saxumVermes/nginx_mansite/pkg/nginx"
)

var info = nginx.Info{
	SitesAvailablePath: "/etc/nginx/sites-available/",
	SitesEnabledPath:   "/etc/nginx/sites-enabled/",
}

func init() {
	enabled, err := ioutil.ReadDir(info.SitesEnabledPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	for _, f := range enabled {
		info.SitesEnabled = append(info.SitesEnabled, f.Name())
	}

	available, err := ioutil.ReadDir(info.SitesAvailablePath)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%v", err)
		os.Exit(1)
	}
	for _, f := range available {
		info.SitesAvailable = append(info.SitesAvailable, f.Name())
	}
}

func main() {
	if len(os.Args) < 2 {
		help("default")
	}

	configCmd := flag.NewFlagSet("config", flag.ExitOnError)
	siteCmd := flag.NewFlagSet("site", flag.ExitOnError)

	switch os.Args[1] {
	case "config":
		c := nginx.Config{}
		configName := configCmd.String("n", "", "Config name")

		if len(os.Args) < 3 {
			help("config")
		}

		switch os.Args[2] {
		case "create":
			configType := configCmd.String("t", "", "Config type: default, drupal")
			configCmd.Parse(os.Args[3:])
			if *configName != "" && *configType != "" {
				c.CreateConfig(&info, *configName, *configType)
			} else {
				configCmd.PrintDefaults()
			}

		case "edit":
			configCmd.Parse(os.Args[3:])
			if *configName != "" {
				c.EditConfig(&info, *configName)
			} else {
				configCmd.PrintDefaults()
			}

		case "delete":
			configCmd.Parse(os.Args[3:])
			if *configName != "" {
				c.DeleteConfig(&info, *configName)
			} else {
				configCmd.PrintDefaults()
			}
		default:
			help("config")
		}

	case "site":
		siteEn := siteCmd.String("e", "", "Enable site")
		siteDis := siteCmd.String("d", "", "Disable site")
		listSiteOf := siteCmd.String("l", "", "List sites: avaliable|enabled")
		siteCmd.Parse(os.Args[2:])

		if len(os.Args) < 3 {
			siteCmd.PrintDefaults()
		}

		if strings.TrimSpace(*listSiteOf) != "" {
			info.ListSites(*listSiteOf)
		}

		if *siteEn != "" {
			info.EnableSite(*siteEn)
		}

		if *siteDis != "" {
			info.DisableSite(*siteDis)
		}

	default:
		help("default")
	}
}

func help(cmd string) {
	switch cmd {
	case "config":
		fmt.Println("\nUsage: nginx-modsite config <params> <flags>")
		fmt.Println("\nParams:\n\tcreate\n\tedit\n\tdelete")

	case "default":
		fmt.Println("\nUsage: nginx-modsite <command> <params> <flags>")
		fmt.Println("\nCommands:\n\tconfig\tcreate, edit, delete configs\n\tsite\tenable or disable sites\n\nType <command> --help for more information.")
	}
	os.Exit(2)
}
