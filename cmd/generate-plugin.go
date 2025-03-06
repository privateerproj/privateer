package cmd

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/revanite-io/sci/pkg/layer2"
)

type CatalogData struct {
	layer2.Catalog
	ServiceName string
	TestSuites  map[string][]string
}

var TemplatesDir string
var SourcePath string
var OutputDir string

// versionCmd represents the version command
var genPluginCmd = &cobra.Command{
	Use:   "generate-plugin",
	Short: "Generate a new plugin",
	Run: func(cmd *cobra.Command, args []string) {
		generatePlugin()
	},
}

func init() {
	genPluginCmd.PersistentFlags().StringP("source-path", "p", "", "The source file to generate the plugin from.")
	genPluginCmd.PersistentFlags().StringP("local-templates", "", "", "Path to a directory to use instead of downloading the latest templates.")
	genPluginCmd.PersistentFlags().StringP("service-name", "n", "", "The name of the service (e.g. 'ECS, AKS, GCS').")
	genPluginCmd.PersistentFlags().StringP("output-dir", "o", "generated-plugin/", "Pathname for the generated plugin.")

	viper.BindPFlag("source-path", genPluginCmd.PersistentFlags().Lookup("source-path"))
	viper.BindPFlag("local-templates", genPluginCmd.PersistentFlags().Lookup("local-templates"))
	viper.BindPFlag("service-name", genPluginCmd.PersistentFlags().Lookup("service-name"))
	viper.BindPFlag("output-dir", genPluginCmd.PersistentFlags().Lookup("output-dir"))

	rootCmd.AddCommand(genPluginCmd)
}

func generatePlugin() {
	err := setupTemplatingEnvironment()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	data, err := readData()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	data.ServiceName = viper.GetString("service-name")
	if data.ServiceName == "" {
		logger.Error("--service-name is required to generate a plugin.")
		return
	}

	err = filepath.Walk(TemplatesDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				err = generateFileFromTemplate(data, path, OutputDir)
				if err != nil {
					logger.Error(fmt.Sprintf("Failed while writing in dir '%s': %s", OutputDir, err))
				}
			} else if info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		},
	)
	if err != nil {
		logger.Error("Error walking through templates directory: %s", err)
	}
}

func setupTemplatingEnvironment() error {
	SourcePath = viper.GetString("source-path")
	if SourcePath == "" {
		return fmt.Errorf("--source-path is required to generate a plugin from a control set from local file or URL")
	}

	if viper.GetString("local-templates") != "" {
		TemplatesDir = viper.GetString("local-templates")
	} else {
		TemplatesDir = filepath.Join(os.TempDir(), "privateer-templates")
		setupTemplatesDir()
	}

	OutputDir = viper.GetString("output-dir")
	logger.Trace(fmt.Sprintf("Generated plugin will be stored in this directory: %s", OutputDir))

	return os.MkdirAll(OutputDir, os.ModePerm)
}

func setupTemplatesDir() error {
	// Remove any old templates
	err := os.RemoveAll(TemplatesDir)
	if err != nil {
		logger.Error("Failed to remove templates directory: %s", err)

	}

	// Pull latest templates from git
	logger.Trace(fmt.Sprintf("Cloning templates repo to: %s", TemplatesDir))
	_, err = git.PlainClone(TemplatesDir, false, &git.CloneOptions{
		URL:      "https://github.com/privateerproj/plugin-generator-templates.git",
		Progress: os.Stdout,
	})
	return err
}

func generateFileFromTemplate(data CatalogData, templatePath, OutputDir string) error {
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("error reading template file %s: %w", templatePath, err)
	}

	tmpl, err := template.New("plugin").Funcs(template.FuncMap{
		"as_text": func(s string) template.HTML {
			s = strings.TrimSpace(strings.ReplaceAll(s, "\n", " "))
			return template.HTML(s)
		},
		"as_id": func(s string) string {
			return strings.TrimSpace(
				strings.ReplaceAll(
					strings.ReplaceAll(s, ".", "_"), "-", "_"))
		},
	}).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("error parsing template file %s: %w", templatePath, err)
	}

	relativePath, err := filepath.Rel(TemplatesDir, templatePath)
	if err != nil {
		return err
	}

	outputPath := filepath.Join(OutputDir, strings.TrimSuffix(relativePath, ".txt"))

	err = os.MkdirAll(filepath.Dir(outputPath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directories for %s: %w", outputPath, err)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file %s: %w", outputPath, err)
	}

	defer outputFile.Close()

	err = tmpl.Execute(outputFile, data)
	if err != nil {
		return fmt.Errorf("error executing template for file %s: %w", outputPath, err)
	}

	return nil
}

func readData() (data CatalogData, err error) {
	err = data.LoadControlFamiliesFile(SourcePath)
	if err != nil {
		return
	}

	data.TestSuites = make(map[string][]string)

	for i, family := range data.ControlFamilies {
		for j := range family.Controls {
			for _, testReq := range data.ControlFamilies[i].Controls[j].Requirements {
				// Add the test ID to the TestSuites map for each TLP level
				for _, tlpLevel := range testReq.Applicability {
					if data.TestSuites[tlpLevel] == nil {
						data.TestSuites[tlpLevel] = []string{}
					}
					data.TestSuites[tlpLevel] = append(data.TestSuites[tlpLevel], testReq.ID)
				}
			}
		}
	}
	return
}
