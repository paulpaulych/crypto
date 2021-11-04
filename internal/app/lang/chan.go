package lang

func HandleChan(errs chan error, onErr func(error)) chan<- error {
	go func() {
		for e := range errs {
			onErr(e)
		}
	}()
	return errs
}