module github.com/privateerproj/privateer

go 1.14

require (
	github.com/hashicorp/go-plugin v1.4.5
	github.com/privateerproj/privateer-sdk v0.0.1-rc // Made this pre-release only to allow for testing brew tap
	github.com/spf13/cobra v1.6.1
	github.com/spf13/viper v1.15.0
	golang.org/x/net v0.7.0 // indirect
)

// For SDK Development Only
// replace github.com/privateerproj/privateer-sdk => ../privateer-sdk
