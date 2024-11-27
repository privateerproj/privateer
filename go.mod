module github.com/privateerproj/privateer

go 1.14

require (
	github.com/go-git/go-git/v5 v5.12.0
	github.com/hashicorp/go-hclog v1.2.0
	github.com/hashicorp/go-plugin v1.4.10
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/privateerproj/privateer-sdk v0.0.7
	github.com/spf13/cobra v1.6.1
	github.com/spf13/viper v1.15.0
	gopkg.in/yaml.v2 v2.4.0
)

// For SDK Development Only
replace github.com/privateerproj/privateer-sdk => ../privateer-sdk
