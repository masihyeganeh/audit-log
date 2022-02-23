package clickhouse

import (
	"context"
)

func (r *repository) AddUser(ctx context.Context, username, hashedPassword, salt string, hasReadAccess, hasWriteAccess bool) error {
	query := `
INSERT INTO users (username, hashed_password, salt, has_read_access, has_write_access)
VALUES ($1,$2,$3,$4,$5);
`
	_, err := r.database.ExecContext(ctx, query, username, hashedPassword, salt, hasReadAccess, hasWriteAccess)
	return err
}

func (r *repository) FindUser(ctx context.Context, username string) (string, string, bool, bool, error) {
	query := `
SELECT hashed_password, salt, has_read_access, has_write_access FROM users
WHERE username = $1;
`
	row := r.database.QueryRowContext(ctx, query, username)

	err := row.Err()
	if err != nil {
		return "", "", false, false, err
	}

	var hashedPassword string
	var salt string
	var hasReadAccess bool
	var hasWriteAccess bool

	err = row.Scan(&hashedPassword, &salt, &hasReadAccess, &hasWriteAccess)
	if err != nil {
		return "", "", false, false, err
	}

	return hashedPassword, salt, hasReadAccess, hasWriteAccess, nil
}
