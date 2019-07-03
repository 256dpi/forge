package forge

// Service manages multiple long-running tasks.
type Service struct {
	Manager
	Terminator
	Reporter
}

// Run will run the specified task continuously until the tasks return or
// service has been stopped or killed.
func (s *Service) Run(n int, task func() error, finalizer func()) {
	s.Manager.Run(n, func() {
		s.Repeat(task)
	}, finalizer)
}
