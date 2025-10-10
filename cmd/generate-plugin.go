package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/ossf/gemara/layer2"
	sdkutils "github.com/privateerproj/privateer-sdk/utils"
)

type CatalogData struct {
	layer2.Catalog
	ServiceName             string
	Requirements            []string
	ApplicabilityCategories []string
	StrippedName            string
}

var (
	TemplatesDir string
	SourcePath   string
	OutputDir    string
	ServiceName  string

	// versionCmd represents the version command
	genPluginCmd = &cobra.Command{
		Use:   "generate-plugin",
		Short: "Generate a new plugin",
		Run: func(cmd *cobra.Command, args []string) {
			generatePlugin()
		},
	}
)

func init() {
	genPluginCmd.PersistentFlags().StringP("source-path", "p", "", "The source file to generate the plugin from.")
	genPluginCmd.PersistentFlags().StringP("local-templates", "", "", "Path to a directory to use instead of downloading the latest templates.")
	genPluginCmd.PersistentFlags().StringP("service-name", "n", "", "The name of the service (e.g. 'ECS, AKS, GCS').")
	genPluginCmd.PersistentFlags().StringP("output-dir", "o", "generated-plugin/", "Pathname for the generated plugin.")

	_ = viper.BindPFlag("source-path", genPluginCmd.PersistentFlags().Lookup("source-path"))
	_ = viper.BindPFlag("local-templates", genPluginCmd.PersistentFlags().Lookup("local-templates"))
	_ = viper.BindPFlag("service-name", genPluginCmd.PersistentFlags().Lookup("service-name"))
	_ = viper.BindPFlag("output-dir", genPluginCmd.PersistentFlags().Lookup("output-dir"))

	rootCmd.AddCommand(genPluginCmd)
}

func generatePlugin() {
	err := setupTemplatingEnvironment()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	data := CatalogData{}
	data.ServiceName = ServiceName

	err = data.LoadFile("file://" + SourcePath)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = data.getAssessmentRequirements()
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

	err = writeCatalogFile(&data.Catalog)
	if err != nil {
		logger.Error("Failed to write catalog to file: %s", err)
	}
}

func setupTemplatingEnvironment() error {
	SourcePath = viper.GetString("source-path")
	if SourcePath == "" {
		return fmt.Errorf("--source-path is required to generate a plugin from a control set from local file or URL")
	}

	ServiceName = viper.GetString("service-name")
	if ServiceName == "" {
		return fmt.Errorf("--service-name is required to generate a plugin.")
	}

	if viper.GetString("local-templates") != "" {
		TemplatesDir = viper.GetString("local-templates")
	} else {
		TemplatesDir = filepath.Join(os.TempDir(), "privateer-templates")
		err := setupTemplatesDir()
		if err != nil {
			return fmt.Errorf("error setting up templates directory: %w", err)
		}
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

	// Determine relative path from templates dir so we can preserve subdirs in output
	relativePath, err := filepath.Rel(TemplatesDir, templatePath)
	if err != nil {
		return fmt.Errorf("error calculating relative path for %s: %w", templatePath, err)
	}

	// If the template is not a text template, copy it over as-is (preserve mode)
	if filepath.Ext(templatePath) != ".txt" {
		return copyNonTemplateFile(templatePath, filepath.Join(OutputDir, relativePath))
	}

	tmpl, err := template.New("plugin").Funcs(template.FuncMap{
		"as_text": func(in string) template.HTML {
			return template.HTML(
				strings.TrimSpace(
					strings.ReplaceAll(in, "\n", " ")))
		},
		"default": func(in string, out string) string {
			if in != "" {
				return in
			}
			return out
		},
		"snake_case":     snakeCase,
		"simplifiedName": simplifiedName,
	}).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("error parsing template file %s: %w", templatePath, err)
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

	defer func() {
		err := outputFile.Close()
		if err != nil {
			logger.Error("error closing output file %s: %w", outputPath, err)
		}
	}()

	err = tmpl.Execute(outputFile, data)
	if err != nil {
		return fmt.Errorf("error executing template for file %s: %w", outputPath, err)
	}

	return nil
}

func (c *CatalogData) getAssessmentRequirements() error {
	for _, family := range c.ControlFamilies {
		for _, control := range family.Controls {
			for _, requirement := range control.AssessmentRequirements {
				c.Requirements = append(c.Requirements, requirement.Id)
				// Add applicability categories if unique
				for _, a := range requirement.Applicability {
					if !sdkutils.StringSliceContains(c.ApplicabilityCategories, a) {
						c.ApplicabilityCategories = append(c.ApplicabilityCategories, a)
					}
				}
			}
		}
	}
	if len(c.Requirements) == 0 {
		return errors.New("no requirements retrieved from catalog")
	}
	return nil
}

func writeCatalogFile(catalog *layer2.Catalog) error {
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(2) // this is the line that sets the indentation
	err := yamlEncoder.Encode(catalog)
	if err != nil {
		return fmt.Errorf("error marshaling YAML: %w", err)
	}

	dirPath := filepath.Join(OutputDir, "data", simplifiedName(catalog.Metadata.Id, catalog.Metadata.Version))
	filePath := filepath.Join(dirPath, "catalog.yaml")

	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directories for %s: %w", filePath, err)
	}

	if err := os.WriteFile(filePath, b.Bytes(), 0644); err != nil {
		return fmt.Errorf("error writing YAML file: %w", err)
	}

	return nil
}

func snakeCase(in string) string {
	return strings.TrimSpace(
		strings.ReplaceAll(
			strings.ReplaceAll(in, ".", "_"), "-", "_"))
}

func simplifiedName(catalogId string, catalogVersion string) string {
	return fmt.Sprintf("%s_%s", snakeCase(catalogId), snakeCase(catalogVersion))
}

func copyNonTemplateFile(templatePath, relativePath string) error {
	outputPath := filepath.Join(OutputDir, relativePath)
	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		return fmt.Errorf("error creating directories for %s: %w", outputPath, err)
	}

	// Copy file contents
	srcFile, err := os.Open(templatePath)
	if err != nil {
		return fmt.Errorf("error opening source file %s: %w", templatePath, err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating destination file %s: %w", outputPath, err)
	}
	defer func() {
		_ = dstFile.Close()
	}()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("error copying file to %s: %w", outputPath, err)
	}

	// Try to preserve file mode from source
	if fi, err := os.Stat(templatePath); err == nil {
		_ = os.Chmod(outputPath, fi.Mode())
	}

	return nil
}
