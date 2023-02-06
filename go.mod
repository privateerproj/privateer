module github.com/privateerproj/privateer

go 1.14

require (
	github.com/hashicorp/go-plugin v1.4.5
	github.com/privateerproj/privateer-sdk v0.0.0 // no release has been made yet
	golang.org/x/net v0.5.0 // indirect
)

// For SDK Development Only
replace github.com/privateerproj/privateer-sdk => ../privateer-sdk
