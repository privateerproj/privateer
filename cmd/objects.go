package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	hcplugin "github.com/hashicorp/go-plugin"
	"github.com/spf13/viper"
)

// RaidError retains an error object and the name of the pack that generated it
type RaidError struct {
	Raid string
	Err  error
}

// RaidErrors holds a list of errors and an Error() method
// so it adheres to the standard Error interface
type RaidErrors struct {
	Errors []RaidError
}

func (e *RaidErrors) Error() string {
	return fmt.Sprintf("Service Pack Errors: %v", e.Errors)
}

type RaidPkg struct {
	Name          string
	Path          string
	ServiceTarget string
	Command       *exec.Cmd
	Result        string

	Available  bool
	Requested  bool
	Successful bool
	Error      error
}

func (p *RaidPkg) getBinary() (binaryName string, err error) {
	p.Name = filepath.Base(strings.ToLower(p.Name)) // in some cases a filepath may arrive here instead of the base name; overwrite if so
	if runtime.GOOS == "windows" && !strings.HasSuffix(p.Name, ".exe") {
		p.Name = fmt.Sprintf("%s.exe", p.Name)
	}
	plugins, _ := hcplugin.Discover(p.Name, viper.GetString("binaries-path"))
	if len(plugins) != 1 {
		err = fmt.Errorf("failed to locate requested plugin '%s' at path '%s'", p.Name, viper.GetString("binaries-path"))
		return
	}
	binaryName = plugins[0]

	return
}

func (p *RaidPkg) queueCmd() {
	cmd := exec.Command(p.Path)
	flags := []string{
		fmt.Sprintf("--config=%s", viper.GetString("config")),
		fmt.Sprintf("--loglevel=%s", viper.GetString("loglevel")),
		fmt.Sprintf("--service=%s", p.ServiceTarget),
	}
	for _, flag := range flags {
		cmd.Args = append(cmd.Args, flag)
		p.Command = cmd
	}
}

func NewRaidPkg(pluginName string, serviceName string) *RaidPkg {
	plugin := &RaidPkg{
		Name: pluginName,
	}
	path, err := plugin.getBinary()
	if err != nil {
		plugin.Error = err
	}
	plugin.Path = path
	plugin.ServiceTarget = serviceName
	plugin.queueCmd()
	return plugin
}
