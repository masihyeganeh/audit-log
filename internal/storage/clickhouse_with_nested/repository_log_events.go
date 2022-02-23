package clickhouse

import (
	"context"
	"github.com/masihyeganeh/audit-log/internal/service"
	logger "log"
)

func (r *repository) LogEvents(ctx context.Context, events []service.Event) error {
	query := `
INSERT INTO events_with_nested (event_time, event_type, common_field_1, common_field_2, fields.Key, fields.Value)
VALUES ($1,$2,$3,$4,$5,$6);
`

	conn, err := r.database.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	statement, err := tx.Prepare(query)
	if err != nil {
		return err
	}

	for _, evt := range events {
		eventToStore := event(evt)
		eventKeys, eventValues := extractKeysAndValues(eventToStore.Fields)
		_, err = statement.Exec(
			eventToStore.EventTime, eventToStore.EventType, eventToStore.CommonField1, eventToStore.CommonField2, eventKeys, eventValues,
		)
		if err != nil {
			logger.Printf("Clickhouse insert failed for event EventType=%q CommonField1=%q CommonField2=%q Fields=%v : %v\n", eventToStore.EventType, eventToStore.CommonField1, eventToStore.CommonField2, eventToStore.Fields, err)
		}
	}

	return tx.Commit()
}
