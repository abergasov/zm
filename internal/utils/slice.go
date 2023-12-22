package utils

func StringsFromObjectSlice[T any](src []T, extractor func(T) string) []string {
	result := make([]string, 0, len(src))
	for key := range src {
		result = append(result, extractor(src[key]))
	}
	return result
}
