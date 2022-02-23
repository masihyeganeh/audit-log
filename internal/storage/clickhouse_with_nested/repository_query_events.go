package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/masihyeganeh/audit-log/internal/service"
	"github.com/pkg/errors"
)

var nonQueryableFields = []string{"event_time", "event_type", "fields"}
var commonFields = []string{"common_field_1", "common_field_2"}

func (r *repository) QueryEvents(ctx context.Context, eventType string, filters map[string]string) ([]service.Event, error) {
	for _, nonQueryableField := range nonQueryableFields {
		delete(filters, nonQueryableField)
	}

	where := "\"event_type\" = $1"
	whereValues := make([]interface{}, 0)
	whereValues = append(whereValues, eventType)

	paramIndex := 2

	for _, commonField := range commonFields {
		if value, exists := filters[commonField]; exists {
			where += fmt.Sprintf(" AND %q = $%d", commonField, paramIndex)
			paramIndex++
			whereValues = append(whereValues, value)
			delete(filters, commonField)
		}
	}

	if len(filters) > 0 {
		for eventKey, eventValue := range filters {
			where += fmt.Sprintf(" AND fields.Value[indexOf(fields.Key, $%d)] = $%d", paramIndex, paramIndex+1)
			paramIndex += 2
			whereValues = append(whereValues, eventKey)
			whereValues = append(whereValues, eventValue)
		}
	}

	query := `
SELECT event_time, event_type, common_field_1, common_field_2, fields.Key, fields.Value
FROM events_with_nested
WHERE %s
LIMIT 1000
`

	query = fmt.Sprintf(query, where)

	results := make([]event, 0)
	rows, err := r.database.QueryContext(ctx, query, whereValues...)
	if err != nil {
		if err == sql.ErrNoRows {
			return []service.Event{}, nil
		}
		return nil, err
	}

	var eventFromDB event
	for rows.Next() {
		var keys []string
		var values []string
		if err = rows.Scan(&eventFromDB.EventTime, &eventFromDB.EventType, &eventFromDB.CommonField1, &eventFromDB.CommonField2, &keys, &values); err != nil {
			return nil, err
		}

		eventFromDB.Fields, err = createMapFromKeysAndValues(keys, values)
		if err != nil {
			return nil, errors.Wrap(err, "could not parse variable fields")
		}

		results = append(results, eventFromDB)
	}

	result := logEvents{Events: results}

	return result.toServiceModel(), nil
}
