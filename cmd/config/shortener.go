package config

type ShortenerConfig struct {
	baseURL    string
	hashLength int
}

func newShortenerConfig() ShortenerConfig {
	return ShortenerConfig{
		baseURL:    getString("SHORTENED_BASE_URL"),
		hashLength: getInt("SHORTENED_URL_HASH_LENGTH"),
	}
}

func (sc ShortenerConfig) GetBaseURL() string {
	return sc.baseURL
}

func (sc ShortenerConfig) GetHashLength() int {
	return sc.hashLength
}
