package utils

import "sync"

// Objek auto increment
type AutoInc struct {
	sync.Mutex // ensures autoInc is goroutine-safe
	Id         int
}

// Kembalikan nilai ID saat ini dan increment untuk ID selanjutnya
func (a *AutoInc) ID() (id int) {
	a.Lock()
	defer a.Unlock()

	id = a.Id
	a.Id++
	return
}
