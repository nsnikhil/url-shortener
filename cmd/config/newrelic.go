package config

type NewRelicConfig struct {
	appName    string
	licenseKey string
}

func (nrc NewRelicConfig) GetAppName() string {
	return nrc.appName
}

func (nrc NewRelicConfig) GetLicenseKey() string {
	return nrc.licenseKey
}

func newNewRelicConfig() NewRelicConfig {
	return NewRelicConfig{
		appName:    getString("NEW_RELIC_APP_NAME"),
		licenseKey: getString("NEW_RELIC_LICENSE_KEY"),
	}
}
