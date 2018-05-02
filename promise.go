package future

type Promise interface {
	OnSuccess()
	OnFailure()
}
