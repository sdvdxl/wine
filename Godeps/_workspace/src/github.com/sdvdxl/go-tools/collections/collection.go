package collections


type IllegalArgumentError struct {
	ErrorMessage string
}

func (e *IllegalArgumentError) Error() string {
	return e.ErrorMessage
}
