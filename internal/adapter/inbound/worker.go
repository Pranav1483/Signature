package inbound

func (h *Handler) StartWorker(done chan any) {
	for {
		select {
		case <-done:
			return
		case f, ok := <-h.worker:
			if !ok {
				return
			}
			go func() {
				h.pool <- struct{}{}
				defer func() {
					<-h.pool
				}()
				f()
			}()
		}
	}
}
