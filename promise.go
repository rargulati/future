package future

// A thennable function
type Promise interface {
	OnSuccess()
	OnFailure()
}
