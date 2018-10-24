package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/saxumVermes/nginx_mansite/src/nginx"
)

//AvailablePath points to nginx directory of available sites.
var AvailablePath = "/etc/nginx/sites-available/"

//EnabledPath points to nginx directory of enabled sites.
var EnabledPath = "/etc/nginx/sites-enabled/"

//TemplatePath is a runtime variable, contains absolute path to templates.
var TemplatePath string

var site = nginx.Site{
	AvailablePath: AvailablePath,
	EnabledPath:   EnabledPath,
}

func init() {
	enabled, err := ioutil.ReadDir(site.EnabledPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	for _, f := range enabled {
		site.Enabled = append(site.Enabled, f.Name())
	}

	available, err := ioutil.ReadDir(site.AvailablePath)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%v", err)
		os.Exit(1)
	}
	for _, f := range available {
		site.Available = append(site.Available, f.Name())
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
		c.TemplatePath = TemplatePath
		configName := configCmd.String("n", "", "Config name")

		if len(os.Args) < 3 {
			help("config")
		}

		switch os.Args[2] {
		case "create":
			configType := configCmd.String("t", "", "Config type: default, drupal")
			configCmd.Parse(os.Args[3:])
			if *configName != "" && *configType != "" {
				c.Create(&site, *configName, *configType)
			} else {
				configCmd.PrintDefaults()
			}

		case "edit":
			configCmd.Parse(os.Args[3:])
			if *configName != "" {
				c.Edit(&site, *configName)
			} else {
				configCmd.PrintDefaults()
			}

		case "delete":
			configs := os.Args[3:]
			if len(configs) > 1 {
				for _, conf := range configs {
					c.Delete(&site, conf)
				}
			} else {
				fmt.Fprintln(os.Stderr, "Invalid arguments. List configs separated by space.")
				os.Exit(1)
			}
		default:
			help("config")
		}

	case "site":
		en := siteCmd.String("e", "", "Enable site")
		dis := siteCmd.String("d", "", "Disable site")
		list := siteCmd.String("l", "", "List sites: avaliable|enabled")
		siteCmd.Parse(os.Args[2:])

		if len(os.Args) < 3 {
			siteCmd.PrintDefaults()
		}

		if strings.TrimSpace(*list) != "" {
			site.List(*list)
		}

		if *en != "" {
			site.Enable(*en)
		}

		if *dis != "" {
			site.Disable(*dis)
		}

	default:
		help("default")
	}
}

func help(cmd string) {
	switch cmd {
	case "config":
		fmt.Println("\nUsage: nginx-man config <params> <flags>")
		fmt.Println("\nParams:\n\tcreate\n\tedit\n\tdelete")

	case "default":
		fmt.Println("\nUsage: nginx-man <command> <params> <flags>")
		fmt.Println("\nCommands:\n\tconfig\tcreate, edit, delete configs\n\tsite\tenable or disable sites\n\nType <command> --help for more information.")
	}
	os.Exit(2)
}
