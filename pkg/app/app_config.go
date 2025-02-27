package app

import (
	"bytes"
	"capten/pkg/config"
	"capten/pkg/types"
	"html/template"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	folderPrmission os.FileMode = 0755
	filePrmission   os.FileMode = 0644
)

func GetClusterGlobalValues(valuesFilePath string) (map[string]interface{}, error) {
	var values map[string]interface{}
	data, err := os.ReadFile(valuesFilePath)
	if err != nil {
		return values, errors.WithMessagef(err, "failed to read cluster values file, %s", valuesFilePath)
	}

	err = yaml.Unmarshal(data, &values)
	if err != nil {
		return values, errors.WithMessagef(err, "failed to unmarshal cluster values file, %s", valuesFilePath)
	}
	return values, nil
}

func GetApps(appListFilePath string) ([]string, error) {
	var values types.AppList
	data, err := os.ReadFile(appListFilePath)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to read app group file, %s", appListFilePath)
	}

	err = yaml.Unmarshal(data, &values)
	if err != nil {
		return nil, errors.WithMessagef(err, "failed to unmarshal app group file, %s", appListFilePath)
	}
	return values.Apps, err
}

func GetAppConfig(appConfigFilePath string, globalValues map[string]interface{}) (types.AppConfig, error) {
	var values types.AppConfig
	data, err := os.ReadFile(appConfigFilePath)
	if err != nil {
		return values, errors.WithMessagef(err, "failed to read app config file, %s", appConfigFilePath)
	}

	transformedData, err := executeAppConfigTemplate(data, globalValues)
	if err != nil {
		return values, errors.WithMessagef(err, "failed to transform app config file, %s", appConfigFilePath)
	}

	err = yaml.Unmarshal(transformedData, &values)
	if err != nil {
		return values, errors.WithMessagef(err, "failed to unmarshal app config file, %s", appConfigFilePath)
	}
	return values, err
}

func WriteAppConfig(captenConfig config.CaptenConfig, appConfig types.AppConfig) error {
	err := os.MkdirAll(captenConfig.PrepareDirPath(captenConfig.AppsTempDirPath), folderPrmission)
	if err != nil {
		return errors.WithMessagef(err, "failed to create directory %s", captenConfig.AppsTempDirPath)
	}

	data, err := yaml.Marshal(&appConfig)
	if err != nil {
		return errors.WithMessagef(err, "failed to unmarshal %s app config", appConfig.Name)
	}

	err = os.WriteFile(captenConfig.PrepareFilePath(captenConfig.AppsTempDirPath, appConfig.Name+".yaml"), data, filePrmission)
	if err != nil {
		return errors.WithMessagef(err, "failed to write %s app config to file", appConfig.Name)
	}
	return nil
}

func PrepareGlobalVaules(captenConfig config.CaptenConfig) (map[string]interface{}, error) {
	globalValues, err := GetClusterGlobalValues(captenConfig.PrepareFilePath(captenConfig.ConfigDirPath, captenConfig.CaptenGlobalValuesFileName))
	if err != nil {
		return nil, err
	}

	err = generateAppGlobalValuesandAppend(globalValues)
	if err != nil {
		return nil, err
	}
	return globalValues, err
}

func executeAppConfigTemplate(data []byte, values map[string]interface{}) (transformedData []byte, err error) {
	tmpl, err := template.New("templateVal").Parse(string(data))
	if err != nil {
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, values)
	if err != nil {
		return
	}

	transformedData = buf.Bytes()
	return
}
