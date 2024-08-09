package cmd

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// ComponentDefinition represents the structure of the input YAML file.
type ComponentDefinition struct {
	ServiceName        string
	CategoryIDFriendly string
	CategoryID         string    `yaml:"category-id"`
	Title              string    `yaml:"title"`
	Controls           []Control `yaml:"controls"`
}

// Control represents the structure of each control within the YAML file.
type Control struct {
	ServiceName      string
	IDFriendly       string
	ID               string              `yaml:"id"`
	FeatureID        string              `yaml:"feature-id"`
	Title            string              `yaml:"title"`
	Objective        string              `yaml:"objective"`
	NISTCSF          string              `yaml:"nist-csf"`
	MITREAttack      string              `yaml:"mitre-attack"`
	ControlMappings  map[string][]string `yaml:"control-mappings"`
	TestRequirements map[string]string   `yaml:"test-requirements"`
}

var Data ComponentDefinition
var TemplatesDir string
var SourcePath string
var OutputDir string

// versionCmd represents the version command
var genRaidCmd = &cobra.Command{
	Use:   "generate-raid",
	Short: "Generate a new raid",
	Run: func(cmd *cobra.Command, args []string) {
		generateRaid()
	},
}

func init() {
	rootCmd.AddCommand(genRaidCmd)

	genRaidCmd.PersistentFlags().StringP("source-path", "p", "", "The source file to generate the raid from.")
	genRaidCmd.PersistentFlags().StringP("local-templates", "", "", "Path to a directory to use instead of downloading the latest templates.")
	genRaidCmd.PersistentFlags().StringP("service-name", "n", "", "The name of the service (e.g. 'ECS, AKS, GCS').")
	genRaidCmd.PersistentFlags().StringP("output-dir", "o", "generated-raid/", "Pathname for the generated raid.")

	viper.BindPFlag("source-path", genRaidCmd.PersistentFlags().Lookup("source-path"))
	viper.BindPFlag("local-templates", genRaidCmd.PersistentFlags().Lookup("local-templates"))
	viper.BindPFlag("service-name", genRaidCmd.PersistentFlags().Lookup("service-name"))
	viper.BindPFlag("output-dir", genRaidCmd.PersistentFlags().Lookup("output-dir"))
}

func generateRaid() {
	err := setupTemplatingEnvironment()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = filepath.Walk(TemplatesDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				err = generateFileFromTemplate(path, OutputDir)
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
		return fmt.Errorf("--source-path is required to generate a raid from a control set from local file or URL.")
	}

	if viper.GetString("local-templates") != "" {
		TemplatesDir = viper.GetString("local-templates")
	} else {
		TemplatesDir = filepath.Join(os.TempDir(), "privateer-templates")
		setupTemplatesDir()
	}

	OutputDir = viper.GetString("output-dir")
	logger.Trace("Generated raid will be stored in this directory: %s", OutputDir)

	if viper.GetString("service-name") == "" {
		return fmt.Errorf("--service-name is required to generate a raid.")
	}
	Data = readData()
	Data.ServiceName = viper.GetString("service-name")

	return os.MkdirAll(OutputDir, os.ModePerm)
}

func setupTemplatesDir() error {
	// Pull latest templates from git
	err := os.RemoveAll(TemplatesDir)
	if err != nil {
		logger.Error("Failed to remove templates directory: %s", err)
	}

	logger.Trace("Cloning templates repo to: ", TemplatesDir)
	_, err = git.PlainClone(TemplatesDir, false, &git.CloneOptions{
		URL:      "https://github.com/privateerproj/raid-generator-templates.git",
		Progress: os.Stdout,
	})
	return err
}

func generateFileFromTemplate(templatePath, OutputDir string) error {
	tmpl, err := template.ParseFiles(templatePath)
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

	err = tmpl.Execute(outputFile, Data)
	if err != nil {
		return fmt.Errorf("error executing template for file %s: %w", outputPath, err)
	}

	return nil
}

func readData() ComponentDefinition {
	var Data ComponentDefinition
	if strings.HasPrefix(SourcePath, "http") {
		Data = readYAMLURL()
	} else {
		Data = readYAMLFile()
	}
	Data.CategoryIDFriendly = strings.ReplaceAll(Data.CategoryID, ".", "_")
	for i := range Data.Controls {
		Data.Controls[i].IDFriendly = strings.ReplaceAll(Data.Controls[i].ID, ".", "_")
	}
	return Data
}

func readYAMLURL() ComponentDefinition {
	resp, err := http.Get(SourcePath)
	if err != nil {
		logger.Error("Failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("Failed to fetch URL: %v", resp.Status)
	}

	var Data ComponentDefinition
	decoder := yaml.NewDecoder(resp.Body)
	err = decoder.Decode(&Data)
	if err != nil {
		logger.Error("Failed to decode YAML from URL: %v", err)
	}

	return Data
}

func readYAMLFile() ComponentDefinition {
	yamlFile, err := os.ReadFile(SourcePath)
	if err != nil {
		logger.Error("Error reading local source file: %s (%v)", SourcePath, err)
	}

	var Data ComponentDefinition
	err = yaml.Unmarshal(yamlFile, &Data)
	if err != nil {
		logger.Error("Error unmarshalling YAML file: %v", err)
	}

	return Data
}
