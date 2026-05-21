package postgres

import (
	"fmt"
	"log/slog"

	"github.com/foxpy/send-me-the-data/src/idb"
)

func (d *Postgres) CreateSessionToken(token idb.SessionToken) error {
	_, err := d.db.Exec(`
		INSERT INTO
			smtd.session_tokens (session_token, expires_at, admin_id)
		VALUES (
			$1,
			$2,
			(SELECT admin_id FROM smtd.admins WHERE username = $3 FOR SHARE)
		)
	`, token.Token, token.ExpiresAt, token.Username)
	if err != nil {
		return fmt.Errorf("failed to create session token: %w", err)
	}

	return nil
}

func (d *Postgres) GetSessionToken(token string) (*idb.SessionToken, error) {
	st := idb.SessionToken{Token: token}
	err := d.db.QueryRow(`
		SELECT
			st.expires_at, a.username
		FROM smtd.session_tokens st
			INNER JOIN smtd.admins a USING (admin_id)
		WHERE st.session_token = $1
	`, token).Scan(&st.ExpiresAt, &st.Username)
	if err != nil {
		return nil, err
	}

	return &st, nil
}

func (d *Postgres) DeleteSessionToken(token string) error {
	_, err := d.db.Exec(`
		DELETE FROM smtd.session_tokens WHERE session_token = $1
	`, token)
	if err != nil {
		return fmt.Errorf("failed to delete session token: %w", err)
	}

	return nil
}

func (d *Postgres) DeleteOutdatedSessionTokens() error {
	res, err := d.db.Exec(`DELETE FROM smtd.session_tokens WHERE expires_at <= CURRENT_TIMESTAMP`)
	if err != nil {
		return fmt.Errorf("failed to delete outdated session tokens: %w", err)
	}

	n, err := res.RowsAffected()
	if err == nil && n > 0 {
		slog.Info("deleted outdated session tokens", "numtokens", n)
	}

	return nil
}
