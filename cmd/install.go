package cmd

import (
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/privateerproj/privateer-sdk/command"
	"github.com/privateerproj/privateer/internal/install"
	"github.com/privateerproj/privateer/internal/registry"
)

var installCmd *cobra.Command

func init() {
	installCmd = &cobra.Command{
		Use:   "install [plugin-name]",
		Short: "Install a vetted plugin from the registry.",
		Long:  "Resolve the plugin name to registry metadata, then download the plugin binary from the release URL into the binaries path.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sdkInstallCmd := command.GetInstallCmd(
				writer,
				func(name string) (*command.PluginMetadata, error) {
					pd, err := registry.NewClient().GetPluginData(name)
					if err != nil {
						return nil, err
					}
					base := strings.TrimSpace(pd.DownloadURL)
					if base == "" {
						if inferred, ok := install.InferGitHubReleaseBase(pd.Source, pd.Latest); ok {
							base = inferred
						} else {
							return &command.PluginMetadata{Name: pd.Name, DownloadURL: ""}, nil
						}
					} else {
						base = strings.TrimSuffix(base, "/")
					}
					binaryName := filepath.Base(name)
					if binaryName == "" || binaryName == "." {
						binaryName = strings.ReplaceAll(name, "/", "-")
					}
					artifactFilename, err := install.InferArtifactFilename(binaryName)
					if err != nil {
						return nil, err
					}
					fullURL := base + "/" + artifactFilename
					return &command.PluginMetadata{
						Name:        pd.Name,
						DownloadURL: fullURL,
					}, nil
				},
				install.FromURL,
			)
			return sdkInstallCmd.RunE(cmd, args)
		},
	}
	rootCmd.AddCommand(installCmd)
}
