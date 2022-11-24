package config

import (
	"fmt"
)

type Config struct {
	Kenda *kendaDbConfig `yaml:"kenda"`
	Cloud *CloudConfigs  `yaml:"cloud"`
}

type CloudConfigs struct {
	Azure Database `yaml:"cloud"`
	// key: factory id, value: container name
	FactoryContainerMap map[string]string `yaml:"-"`
	// key: container name
	ContainerToken map[string]BlobToken `yaml:"-"`

	FactoryContainerInfo []FactoryContainerInfo `yaml:"factoryContainerInfo"`

	UploadRetryParameters UploadRetryParameters `yaml:"uploadRetryParameters"`
}

func (ccf CloudConfigs) Validate() error {
	for _, containerName := range ccf.FactoryContainerMap {
		if _, ok := ccf.ContainerToken[containerName]; !ok {
			return fmt.Errorf("insufficient cloud configs, container:" + containerName + " token not found")
		}
	}
	return nil
}

func (ccf *CloudConfigs) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := unmarshal(ccf); err != nil {
		return err
	}

	ccf.parseFactoryContainerInfo()

	return nil
}

func (ccf *CloudConfigs) parseFactoryContainerInfo() {
	ccf.FactoryContainerMap = make(map[string]string)
	ccf.ContainerToken = make(map[string]BlobToken)

	for _, info := range ccf.FactoryContainerInfo {
		if _, ok := ccf.FactoryContainerMap[info.FactoryID]; !ok {
			ccf.FactoryContainerMap[info.FactoryID] = info.ContainerName
		}

		if _, ok := ccf.ContainerToken[info.ContainerName]; !ok {
			ccf.ContainerToken[info.ContainerName] = info.BlobToken
		}
	}
}

type FactoryContainerInfo struct {
	FactoryID     string    `yaml:"factoryID"`
	ContainerName string    `yaml:"containerName"`
	BlobToken     BlobToken `yaml:"blobToken"`
}

type UploadRetryParameters struct {
	MaxTimes int `yaml:"maxTimes"`
	// unit: minutes
	Interval             int    `yaml:"interval"`
	MaxConcurrentJobs    int    `yaml:"maxConcurrentJobs"`
	UnsuccessStoragePath string `yaml:"unsuccessStoragePath"`
}

type BlobToken struct {
	AccountName string `yaml:"accountName"`
	SasToken    string `yaml:"sasToken"`
}

type kendaDbConfig map[string]Database

type Database struct {
	Schema   string `yaml:"schema"`
	Name     string `yaml:"name"`
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
}
