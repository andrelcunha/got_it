package config

func GetAcceptedKeys() []string {
	keys := make([]string, 0, len(acceptedKeys))
	for key := range acceptedKeys {
		keys = append(keys, key)
	}
	return keys
}

func IsValidKey(key string) bool {
	_, ok := acceptedKeys[key]
	return ok
}

func GetEssentilFiles() []string {
	return essentialFiles
}
