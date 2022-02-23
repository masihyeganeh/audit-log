package clickhouse

import (
	"context"
	"database/sql"
)

func (r *repository) setup(ctx context.Context) (error, bool) {
	tx, err := r.database.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err, false
	}

	// events
	tx.Exec(`
CREATE TABLE IF NOT EXISTS events_with_nested
(
    event_time DateTime,
    event_type String,
    common_field_1 String,
    common_field_2 String,
    fields Nested (
        Key String,
        Value String
    )
) ENGINE = MergeTree()
PARTITION BY event_type
ORDER BY tuple();
`)
	// users
	tx.Exec(`
CREATE TABLE IF NOT EXISTS users
(
    username String,
    hashed_password String,
    salt String,
    has_read_access Boolean,
    has_write_access Boolean
) ENGINE = MergeTree()
PARTITION BY username
ORDER BY (username);
`)

	err = tx.Commit()
	if err != nil {
		return err, false
	}

	query := `
SELECT count()
FROM users;
`

	var countOfUsers int
	row := r.database.QueryRowContext(ctx, query)

	err = row.Err()
	if err != nil {
		return err, false
	}

	err = row.Scan(&countOfUsers)
	if err != nil {
		return err, false
	}

	return nil, countOfUsers == 0
}
