package utils

type Semaphore chan Empty

func (s Semaphore) aquire(n int) {
	e := Empty{}
	for i := 0; i < n; i++ {
		s <- e
	}
}

// release n resources
func (s Semaphore) release(n int) {
	for i := 0; i < n; i++ {
		<-s
	}
}

/* mutexes */

func (s Semaphore) Lock() {
	s.aquire(1)
}

func (s Semaphore) Unlock() {
	s.release(1)
}

/* signal-wait */

func (s Semaphore) Signal() {
	s.release(1)
}

func (s Semaphore) Wait(n int) {
	s.aquire(n)
}
