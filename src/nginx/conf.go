package nginx

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"
)

//Config contains nginx site configurations.
type Config struct {
	Port         int
	ServerName   string
	Root         string
	Templates    map[string]*template.Template
	TemplatePath string
}

//Edit provides a faster way to edit nginx config in vi.
func (c *Config) Edit(s *Site, name string) {
	file, err := os.OpenFile(s.AvailablePath+strings.TrimSpace(name), os.O_APPEND|os.O_RDWR, 0744)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while editing file %s: %v\n", name, err)
		os.Exit(1)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			os.Exit(1)
		}
	}()

	editor, err := exec.LookPath(strings.TrimSpace("vi"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "vi not found: %v", err)
		os.Exit(1)
	}
	exec.Command(editor, s.AvailablePath+file.Name()).Run()
}

//Create creates new configuration from a template based on the provided answers.
func (c *Config) Create(s *Site, name string, templateType string) {
	c.Templates = make(map[string]*template.Template)
	c.Templates["default"] = template.Must(template.ParseFiles(fmt.Sprintf("%s/templates/default.conf", c.TemplatePath)))
	c.Templates["drupal"] = template.Must(template.ParseFiles(fmt.Sprintf("%s/templates/drupal.conf", c.TemplatePath)))

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("You are about to create a new configuration. Give a port number: ")
	p, _, err := reader.ReadLine()
	c.Port, _ = strconv.Atoi(string(p))

	reader.Reset(os.Stdin)
	fmt.Print("Give a server name: ")
	c.ServerName, err = reader.ReadString('\n')

	reader.Reset(os.Stdin)
	fmt.Print("Server root: ")
	c.Root, err = reader.ReadString('\n')

	if err != nil {
		fmt.Fprintf(os.Stdout, "%v", err)
		os.Exit(1)
	}

	file, err := os.OpenFile(s.AvailablePath+strings.TrimSpace(name), os.O_CREATE|os.O_RDWR, 0744)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not create file %s: %v\n", name, err)
		os.Exit(1)
	}
	defer file.Close()
	if t, ok := c.Templates[templateType]; ok {
		t.Execute(file, c)
	} else {
		fmt.Fprintln(os.Stderr, "Such configuration does not exist!")
		os.Exit(1)
	}
	fmt.Println("Config successfully created!")
}

//Delete removes a specified configuration form sites available directory.
func (c *Config) Delete(s *Site, name string) {
	var conf string
	for _, s := range s.Available {
		if s == strings.TrimSpace(name) {
			conf = name
		}
	}
	if conf == "" {
		fmt.Fprintln(os.Stderr, "Config not found!")
		os.Exit(1)
	}
	err := os.Remove(s.AvailablePath + conf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to remove %s: %v\n", conf, err)
		os.Exit(1)
	}
	fmt.Println("Config deleted!")
}
