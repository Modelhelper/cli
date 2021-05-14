package app

func Initialize(initializer Initializer) error {
	err := initializer.Initialize()
	return err
}

type Initializer interface {
	Initialize() error
}

func (a *Application) Initialize(init ...Initializer) error {
	for _, i := range init {
		err := i.Initialize()
		if err != nil {
			return err
		}
	}

	return nil
}
