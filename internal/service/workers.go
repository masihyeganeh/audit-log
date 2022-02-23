package service

import (
	"context"
	"log"
	"time"
)

func (s *service) StartBackgroundWorkers(ctx context.Context) {
	for i := 0; i < s.workersCount; i++ {
		go s.spinUpWorker(ctx)
	}
	<-ctx.Done()
}

func (s *service) spinUpWorker(ctx context.Context) {
	buffer := make([]Event, 0, s.eventsBufferSize)
	t := time.NewTicker(s.maxEventsBufferDuration)
	for {
		select {
		case <-t.C:
			if len(buffer) > 0 {
				err := s.datastoreRepository.LogEvents(ctx, buffer)
				if err != nil {
					log.Printf("could not insert events to datastore : %v\n", err)
				}
				buffer = make([]Event, 0, s.eventsBufferSize)
			}
		case event := <-s.channel:
			buffer = append(buffer, event)
			if len(buffer) == s.eventsBufferSize {
				err := s.datastoreRepository.LogEvents(ctx, buffer)
				if err != nil {
					log.Printf("could not insert events to datastore : %v\n", err)
				}
				buffer = make([]Event, 0, s.eventsBufferSize)
				t.Reset(s.maxEventsBufferDuration)
			}
		case <-ctx.Done():
			t.Stop()
			if len(buffer) > 0 {
				err := s.datastoreRepository.LogEvents(ctx, buffer)
				if err != nil {
					log.Printf("could not insert events to datastore : %v\n", err)
				}
			}
			return
		}
	}
}
