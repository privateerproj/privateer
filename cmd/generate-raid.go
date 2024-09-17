package cmd

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type ComponentDefinition struct {
	CategoryIDFriendly string
	ServiceName        string
	Metadata           Metadata  `yaml:"metadata"`
	Controls           []Control `yaml:"controls"`
	Features           []Feature `yaml:"features"`
	Threats            []Threat  `yaml:"threats"`
}

type Control struct {
	IDFriendly       string
	ID               string                 `yaml:"id"`
	Title            string                 `yaml:"title"`
	Objective        string                 `yaml:"objective"`
	ControlFamily    string                 `yaml:"control_family"`
	Threats          []string               `yaml:"threats"`
	NISTCSF          string                 `yaml:"nist_csf"`
	MITREATTACK      string                 `yaml:"mitre_attack"`
	ControlMappings  map[string]interface{} `yaml:"control_mappings"`
	TestRequirements map[string]string      `yaml:"test_requirements"`
}

// Metadata is a struct that represents the metadata.yaml file
type Metadata struct {
	Title              string `yaml:"title"`
	ID                 string `yaml:"id"`
	Description        string `yaml:"description"`
	AssuranceLevel     string `yaml:"assurance_level"`
	ThreatModelAuthor  string `yaml:"threat_model_author"`
	ThreatModelURL     string `yaml:"threat_model_url"`
	RedTeam            string `yaml:"red_team"`
	RedTeamExercizeURL string `yaml:"red_team_exercise_url"`
}

type Feature struct {
	ID          string `yaml:"id"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}

type Threat struct {
	ID          string   `yaml:"id"`
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Features    []string `yaml:"features"`
	MITRE       []string `yaml:"mitre_attack"`
}

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
	data, err := readData()
	if err != nil {
		log.Error(err.Error())
		return
	}
	data.ServiceName = viper.GetString("service-name")
	if data.ServiceName == "" {
		log.Error(fmt.Errorf("--service-name is required to generate a raid."))
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
		return fmt.Errorf("--source-path is required to generate a raid from a control set from local file or URL.")
	}

	if viper.GetString("local-templates") != "" {
		TemplatesDir = viper.GetString("local-templates")
	} else {
		TemplatesDir = filepath.Join(os.TempDir(), "privateer-templates")
		setupTemplatesDir()
	}

	OutputDir = viper.GetString("output-dir")
	logger.Trace(fmt.Sprintf("Generated raid will be stored in this directory: %s", OutputDir))

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
		URL:      "https://github.com/privateerproj/raid-generator-templates.git",
		Progress: os.Stdout,
	})
	return err
}

func generateFileFromTemplate(data ComponentDefinition, templatePath, OutputDir string) error {
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

	err = tmpl.Execute(outputFile, data)
	if err != nil {
		return fmt.Errorf("error executing template for file %s: %w", outputPath, err)
	}

	return nil
}

func readData() (data ComponentDefinition, err error) {
	if strings.HasPrefix(SourcePath, "http") {
		data, err = readYAMLURL()
	} else {
		data, err = readYAMLFile()
	}
	if err != nil {
		return
	}

	data.CategoryIDFriendly = strings.ReplaceAll(data.Metadata.ID, ".", "_")
	for i := range data.Controls {
		data.Controls[i].IDFriendly = strings.ReplaceAll(data.Controls[i].ID, ".", "_")
		// loop over objectives in test_requirements and replace newlines with empty string
		for k, v := range data.Controls[i].TestRequirements {
			data.Controls[i].TestRequirements[k] = strings.Replace(v, "\n", "", -1)
		}
	}
	return
}

func readYAMLURL() (data ComponentDefinition, err error) {
	resp, err := http.Get(SourcePath)
	if err != nil {
		logger.Error("Failed to fetch URL: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error("Failed to fetch URL: %v", resp.Status)
		return
	}

	decoder := yaml.NewDecoder(resp.Body)
	err = decoder.Decode(&data)
	if err != nil {
		logger.Error("Failed to decode YAML from URL: %v", err)
		return
	}

	return
}

func readYAMLFile() (data ComponentDefinition, err error) {
	yamlFile, err := os.ReadFile(SourcePath)
	if err != nil {
		logger.Error(fmt.Sprintf("Error reading local source file: %s (%v)", SourcePath, err))
		return
	}

	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		logger.Error("Error unmarshalling YAML file: %v", err)
		return
	}

	return
}
