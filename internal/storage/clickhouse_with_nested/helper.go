package clickhouse

import "github.com/pkg/errors"

func extractKeysAndValues(m map[string]string) ([]string, []string) {
	keys := make([]string, 0, len(m))
	values := make([]string, 0, len(m))
	for k, v := range m {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}

func createMapFromKeysAndValues(keys []string, values []string) (map[string]string, error) {
	if len(keys) != len(values) {
		return map[string]string{}, errors.New("count of keys and values are not the same")
	}

	result := make(map[string]string, len(keys))

	for i := 0; i < len(keys); i++ {
		result[keys[i]] = values[i]
	}

	return result, nil
}
