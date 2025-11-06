package concurrency

import "library_management/services"

// StartConcurrentReservationWorker runs reservation handling in the background
func StartConcurrentReservationWorker(library *services.Library) {
	go library.StartReservationWorker()
}
