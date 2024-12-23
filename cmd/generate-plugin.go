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

type ControlCatalog struct {
	CategoryIDFriendly string
	ServiceName        string
	TestSuites            map[string][]string

	Metadata Metadata `yaml:"metadata"`

	Controls []Control `yaml:"controls"`
	Features []Feature `yaml:"features"`
	Threats  []Threat  `yaml:"threats"`

	LatestReleaseDetails ReleaseDetails `yaml:"latest_release_details"`
}

// Metadata is a struct that represents the metadata.yaml file
type Metadata struct {
	Title          string           `yaml:"title"`
	ID             string           `yaml:"id"`
	Description    string           `yaml:"description"`
	ReleaseDetails []ReleaseDetails `yaml:"release_details"`
}

type ReleaseDetails struct {
	Version            string         `yaml:"version"`
	AssuranceLevel     string         `yaml:"assurance_level"`
	ThreatModelURL     string         `yaml:"threat_model_url"`
	ThreatModelAuthor  string         `yaml:"threat_model_author"`
	RedTeam            string         `yaml:"red_team"`
	RedTeamExerciseURL string         `yaml:"red_team_exercise_url"`
	ReleaseManager     ReleaseManager `yaml:"release_manager"`
	ChangeLog          []string       `yaml:"change_log"`
}

type ReleaseManager struct {
	Name     string `yaml:"name"`
	GithubId string `yaml:"github_id"`
	Company  string `yaml:"company"`
	Summary  string `yaml:"summary"`
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
	TestRequirements []TestRequirement      `yaml:"test_requirements"`
}

type TestRequirement struct {
	IDFriendly string
	ID         string   `yaml:"id"`
	Text       string   `yaml:"text"`
	TLPLevels  []string `yaml:"tlp_levels"`
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
	rootCmd.AddCommand(genPluginCmd)

	genPluginCmd.PersistentFlags().StringP("source-path", "p", "", "The source file to generate the plugin from.")
	genPluginCmd.PersistentFlags().StringP("local-templates", "", "", "Path to a directory to use instead of downloading the latest templates.")
	genPluginCmd.PersistentFlags().StringP("service-name", "n", "", "The name of the service (e.g. 'ECS, AKS, GCS').")
	genPluginCmd.PersistentFlags().StringP("output-dir", "o", "generated-plugin/", "Pathname for the generated plugin.")

	viper.BindPFlag("source-path", genPluginCmd.PersistentFlags().Lookup("source-path"))
	viper.BindPFlag("local-templates", genPluginCmd.PersistentFlags().Lookup("local-templates"))
	viper.BindPFlag("service-name", genPluginCmd.PersistentFlags().Lookup("service-name"))
	viper.BindPFlag("output-dir", genPluginCmd.PersistentFlags().Lookup("output-dir"))
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

func generateFileFromTemplate(data ControlCatalog, templatePath, OutputDir string) error {
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

func readData() (data ControlCatalog, err error) {
	if strings.HasPrefix(SourcePath, "http") {
		data, err = readYAMLURL()
	} else {
		data, err = readYAMLFile()
	}
	if err != nil {
		return
	}

	data.TestSuites = make(map[string][]string)
	data.CategoryIDFriendly = strings.ReplaceAll(data.Metadata.ID, ".", "_")

	for i := range data.Controls {
		fmt.Println(data.Controls[i].ID)
		data.Controls[i].IDFriendly = strings.ReplaceAll(data.Controls[i].ID, ".", "_")
		// loop over objectives in test_requirements and replace newlines with empty string
		for j, testReq := range data.Controls[i].TestRequirements {
			// Some test requirements have newlines in them, which breaks the template
			data.Controls[i].TestRequirements[j].Text = strings.TrimSpace(strings.ReplaceAll(testReq.Text, "\n", " "))
			// Replace periods with underscores for the friendly ID
			data.Controls[i].TestRequirements[j].IDFriendly = strings.ReplaceAll(testReq.ID, ".", "_")

			// Add the test ID to the TestSuites map for each TLP level
			for _, tlpLevel := range testReq.TLPLevels {
				if data.TestSuites[tlpLevel] == nil {
					data.TestSuites[tlpLevel] = []string{}
				}
				data.TestSuites[tlpLevel] = append(data.TestSuites[tlpLevel], strings.ReplaceAll(testReq.ID, ".", "_"))
			}
		}
	}
	return
}

func readYAMLURL() (data ControlCatalog, err error) {
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

func readYAMLFile() (data ControlCatalog, err error) {
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
