package cmd

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

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
	genRaidCmd.PersistentFlags().StringP("service-name", "n", "", "The name of the service.")
	genRaidCmd.PersistentFlags().StringP("output-dir", "o", "", "The name of the service.")

	viper.BindPFlag("source-path", genRaidCmd.PersistentFlags().Lookup("source-path"))
	viper.BindPFlag("service-name", genRaidCmd.PersistentFlags().Lookup("service-name"))
	viper.BindPFlag("output-dir", genRaidCmd.PersistentFlags().Lookup("output-dir"))

}

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

func generateRaid() {
	// TODO: Pull this from a stable repo, or otherwise get it from the binary somehow
	TemplatesDir = filepath.Join("cmd", "templates")

	SourcePath = viper.GetString("source-path")
	OutputDir = viper.GetString("output-dir")

	Data = readData()
	Data.ServiceName = viper.GetString("service-name")

	err := os.MkdirAll(OutputDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create output directory: %s", err)
	}

	err = filepath.Walk(TemplatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			err = generateFileFromTemplate(path, OutputDir)
			if err != nil {
				log.Fatal(fmt.Sprintf("Failed while writing in dir '%s': %s", OutputDir, err))
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error walking through templates directory: %s", err)
	}
	fmt.Println("Go project directory generated successfully.")
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
		log.Fatalf("Failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to fetch URL: %v", resp.Status)
	}

	var Data ComponentDefinition
	decoder := yaml.NewDecoder(resp.Body)
	err = decoder.Decode(&Data)
	if err != nil {
		log.Fatalf("Failed to decode YAML from URL: %v", err)
	}

	return Data
}

func readYAMLFile() ComponentDefinition {
	yamlFile, err := os.ReadFile(SourcePath)
	if err != nil {
		log.Fatalf("Error reading local source file: %s (%v)", SourcePath, err)
	}

	var Data ComponentDefinition
	err = yaml.Unmarshal(yamlFile, &Data)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML file: %v", err)
	}

	return Data
}
