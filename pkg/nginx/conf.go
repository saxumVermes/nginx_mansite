package nginx

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

type Config struct {
	Port       int
	ServerName string
	Root       string
}

func (c *Config) EditConfig(i *Info, name string) {
	file, err := os.OpenFile(i.SitesAvailablePath+strings.TrimSpace(name), os.O_APPEND|os.O_RDWR, 0744)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not create/edit file %s: %v\n", name, err)
		os.Exit(1)
	}
	defer file.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter an executable editor: ")
	response, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	editor, err := exec.LookPath(strings.TrimSpace(response))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Editor not found: %v", err)
		os.Exit(1)
	}
	exec.Command(editor, i.SitesAvailablePath+file.Name()).Run()
}

func (c *Config) CreateConfig(i *Info, name string, templateType string) {
	temps := make(map[string]*template.Template)
	temps["default"] = template.Must(template.ParseFiles("../../templates/default.conf"))
	temps["drupal"] = template.Must(template.ParseFiles("../../templates/drupal.conf"))

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("You are about to create a new configuration. Give a port number: ")
	p, err := reader.ReadByte()
	c.Port = int(p)

	reader.Reset(os.Stdin)
	fmt.Print("Give a server name: ")
	c.ServerName, err = reader.ReadString('\n')

	reader.Reset(os.Stdin)
	fmt.Print("Server root: ")
	c.Root, err = reader.ReadString('\n')

	if err != nil {
		fmt.Fprintf(os.Stdout, "%v", err)
	}

	file, err := os.OpenFile(i.SitesAvailablePath+strings.TrimSpace(name), os.O_CREATE|os.O_RDWR, 0744)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not create file %s: %v\n", name, err)
		os.Exit(1)
	}
	defer file.Close()
	if t, ok := temps[templateType]; ok {
		t.Execute(file, c)
	} else {
		fmt.Fprintln(os.Stderr, "Such configuration does not exist!")
		os.Exit(1)
	}
}

func (c *Config) DeleteConfig(i *Info, name string) {
	var conf string
	for _, s := range i.SitesAvailable {
		if s == strings.TrimSpace(name) {
			conf = name
		}
	}
	if conf == "" {
		fmt.Fprintln(os.Stderr, "Config not found!")
		os.Exit(1)
	}
	err := os.Remove(i.SitesAvailablePath + conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to remove %s: %v\n", conf, err)
		os.Exit(1)
	}
	fmt.Println("Config deleted!")
}
