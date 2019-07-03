package forge

// Service combines Manager, Terminator and Reporter.
type Service struct {
	Manager
	Terminator
	Reporter
}

// Run will run the specified task continuously until the service has been
// stopped or killed.
func (s *Service) Run(n int, task func() error, finalizer func()) {
	s.Manager.Run(n, func() {
		s.Repeat(task)
	}, finalizer)
}
