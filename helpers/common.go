package helpers

func ResultOrEmpty[T any](input []T) []T {
	if input == nil {
		return []T{}
	}
	return input
}
