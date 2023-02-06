package config

import (
	"log"
	"os"
	"path/filepath"

	sdkConfig "github.com/privateerproj/privateer-sdk/config"
	"github.com/privateerproj/privateer-sdk/config/setter"
	"github.com/privateerproj/privateer-sdk/logging"
)

type varOptions struct {
	VarsFile *string

	AllRaids     *bool    `yaml:"AllRaids"`
	Verbose      *bool    `yaml:"Verbose"`
	BinariesPath string   `yaml:"BinariesPath"`
	Run          []string `yaml:"Run"`
}

// Vars is a stateful object containing the variables required to execute this pack
var Vars varOptions

// Init will set values with the content retrieved from a filepath, env vars, or defaults
func (ctx *varOptions) Init() (err error) {

	if ctx.varsFileIsFound() {
		sdkConfig.GlobalConfig.VarsFile = *ctx.VarsFile
		ctx.decode()
		if err != nil {
			log.Printf("[ERROR] %v", err)
			return
		}
	} else {
		log.Printf("[WARN] No vars file provided, unexpected behavior may occur")
	}
	sdkConfig.GlobalConfig.Init()
	logging.UseLogger("core")
	sdkConfig.GlobalConfig.PrepareOutputDirectory() // TODO: Move this to a better location, currently runs on init (including -h)
	ctx.setEnvAndDefaults()
	return
}

func (ctx *varOptions) varsFileIsFound() bool {
	if ctx.VarsFile == nil {
		defaultFilename := "config.yml"
		ctx.VarsFile = &defaultFilename
	}
	_, err := os.Stat(*ctx.VarsFile)
	return err == nil
}

// decode uses an SDK helper to create a YAML file decoder,
// parse the file to an object, then extracts the values from
// Raid.Kubernetes into this context
func (ctx *varOptions) decode() (err error) {
	configDecoder, file, err := sdkConfig.NewConfigDecoder(*ctx.VarsFile)
	if err != nil {
		return
	}
	err = configDecoder.Decode(&ctx)
	file.Close()
	return err
}

func (ctx *varOptions) setEnvAndDefaults() {
	setter.SetVar(&ctx.BinariesPath, "PRIVATEER_BIN", filepath.Join(sdkConfig.GlobalConfig.InstallDir, "bin"))

	f := false
	if ctx.AllRaids == nil {
		ctx.AllRaids = &f
	}
	if ctx.Verbose == nil {
		ctx.Verbose = &f
	}
}
