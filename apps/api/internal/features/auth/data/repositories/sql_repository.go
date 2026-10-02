package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aisha-platform/aisha/apps/api/internal/features/auth/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User = domain.User
type Session = domain.Session

var (
	ErrInvalidCredentials = domain.ErrInvalidCredentials
	ErrEmailExists        = domain.ErrEmailExists
	ErrInvalidToken       = domain.ErrInvalidToken
)

type PostgresRepository struct{ pool *pgxpool.Pool }

const emailVerificationLifetime = "120 seconds"

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, email, passwordHash, name, verificationHash, activationURL string) (User, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return User{}, fmt.Errorf("begin create user: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var user User
	err = tx.QueryRow(ctx, `INSERT INTO users(email,password_hash,display_name) VALUES($1,$2,$3) RETURNING id,email,password_hash,display_name,status,email_verified_at,created_at`, email, passwordHash, name).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.DisplayName, &user.Status, &user.EmailVerifiedAt, &user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return User{}, ErrEmailExists
		}
		return User{}, fmt.Errorf("insert user: %w", err)
	}
	if _, err = tx.Exec(ctx, `INSERT INTO user_roles(user_id,role_id) SELECT $1,id FROM roles WHERE code='customer'`, user.ID); err != nil {
		return User{}, fmt.Errorf("assign customer role: %w", err)
	}
	if _, err = tx.Exec(ctx, `INSERT INTO email_verification_tokens(user_id,token_hash,expires_at) VALUES($1,$2,CURRENT_TIMESTAMP + $3::interval)`, user.ID, verificationHash, emailVerificationLifetime); err != nil {
		return User{}, fmt.Errorf("store email verification token: %w", err)
	}
	// Keep the aggregate id parameter consistently typed as UUID. Reusing the
	// same placeholder as both UUID and text makes PostgreSQL reject the
	// statement with SQLSTATE 42P08 before the registration can commit.
	if _, err = tx.Exec(ctx, `INSERT INTO outbox_events(event_type,aggregate_type,aggregate_id,payload) VALUES('USER_REGISTERED','user',$1::uuid,jsonb_build_object('userId',$1::uuid,'email',$2::text,'activationUrl',$3::text))`, user.ID, user.Email, activationURL); err != nil {
		return User{}, fmt.Errorf("write registration notification: %w", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return User{}, fmt.Errorf("commit create user: %w", err)
	}
	user.Roles = []string{"customer"}
	return user, nil
}
func (r *PostgresRepository) UserByEmail(ctx context.Context, email string) (User, error) {
	return r.user(ctx, `WHERE lower(u.email)=lower($1)`, email)
}
func (r *PostgresRepository) UserByIdentifier(ctx context.Context, identifier string) (User, error) {
	identifier = strings.TrimSpace(identifier)
	if strings.Contains(identifier, "@") {
		return r.UserByEmail(ctx, identifier)
	}
	phone := normalizedPhone(identifier)
	if phone == "" {
		return User{}, ErrInvalidCredentials
	}
	storedPhone := `regexp_replace(COALESCE(u.phone,''),'[^0-9]','','g')`
	canonicalStoredPhone := `CASE WHEN ` + storedPhone + ` LIKE '00213%' THEN '0'||substring(` + storedPhone + ` FROM 6) WHEN ` + storedPhone + ` LIKE '213%' THEN '0'||substring(` + storedPhone + ` FROM 4) ELSE ` + storedPhone + ` END`
	return r.user(ctx, `WHERE `+canonicalStoredPhone+`=$1`, phone)
}
func (r *PostgresRepository) UserByID(ctx context.Context, id string) (User, error) {
	return r.user(ctx, `WHERE u.id=$1`, id)
}
func (r *PostgresRepository) user(ctx context.Context, where, arg string) (User, error) {
	var user User
	err := r.pool.QueryRow(ctx, `SELECT u.id,u.email,COALESCE(u.phone,''),u.password_hash,u.display_name,u.status,u.email_verified_at,u.created_at,COALESCE(array_agg(r.code) FILTER(WHERE r.code IS NOT NULL),'{}') FROM users u LEFT JOIN user_roles ur ON ur.user_id=u.id LEFT JOIN roles r ON r.id=ur.role_id `+where+` GROUP BY u.id`, arg).Scan(&user.ID, &user.Email, &user.Phone, &user.PasswordHash, &user.DisplayName, &user.Status, &user.EmailVerifiedAt, &user.CreatedAt, &user.Roles)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, ErrInvalidCredentials
	}
	if err != nil {
		return User{}, fmt.Errorf("query user: %w", err)
	}
	return user, nil
}

func normalizedPhone(value string) string {
	var digits strings.Builder
	for _, character := range value {
		if character >= '0' && character <= '9' {
			digits.WriteRune(character)
		}
	}
	phone := digits.String()
	switch {
	case strings.HasPrefix(phone, "00213"):
		return "0" + phone[5:]
	case strings.HasPrefix(phone, "213"):
		return "0" + phone[3:]
	default:
		return phone
	}
}
func (r *PostgresRepository) CreateSession(ctx context.Context, userID, tokenHash, familyID string, expires time.Time) (Session, error) {
	var session Session
	err := r.pool.QueryRow(ctx, `INSERT INTO sessions(user_id,token_hash,family_id,expires_at) VALUES($1,$2,$3,$4) RETURNING id,user_id,family_id,expires_at,revoked_at`, userID, tokenHash, familyID, expires).Scan(&session.ID, &session.UserID, &session.FamilyID, &session.ExpiresAt, &session.RevokedAt)
	if err != nil {
		return Session{}, fmt.Errorf("insert session: %w", err)
	}
	return session, nil
}
func (r *PostgresRepository) RotateSession(ctx context.Context, currentHash, nextHash string, expires time.Time) (Session, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return Session{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var old Session
	err = tx.QueryRow(ctx, `SELECT id,user_id,family_id,expires_at,revoked_at FROM sessions WHERE token_hash=$1 FOR UPDATE`, currentHash).Scan(&old.ID, &old.UserID, &old.FamilyID, &old.ExpiresAt, &old.RevokedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return Session{}, ErrInvalidToken
	}
	if err != nil {
		return Session{}, fmt.Errorf("lock session: %w", err)
	}
	now := time.Now().UTC()
	if old.RevokedAt != nil || !old.ExpiresAt.After(now) {
		_, _ = tx.Exec(ctx, `UPDATE sessions SET revoked_at=COALESCE(revoked_at,$2) WHERE family_id=$1`, old.FamilyID, now)
		_ = tx.Commit(ctx)
		return Session{}, ErrInvalidToken
	}
	var next Session
	err = tx.QueryRow(ctx, `INSERT INTO sessions(user_id,token_hash,family_id,expires_at) VALUES($1,$2,$3,$4) RETURNING id,user_id,family_id,expires_at,revoked_at`, old.UserID, nextHash, old.FamilyID, expires).Scan(&next.ID, &next.UserID, &next.FamilyID, &next.ExpiresAt, &next.RevokedAt)
	if err != nil {
		return Session{}, fmt.Errorf("insert replacement session: %w", err)
	}
	if _, err = tx.Exec(ctx, `UPDATE sessions SET revoked_at=$2,replaced_by_session_id=$3,last_used_at=$2 WHERE id=$1`, old.ID, now, next.ID); err != nil {
		return Session{}, fmt.Errorf("revoke replaced session: %w", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return Session{}, fmt.Errorf("commit rotation: %w", err)
	}
	return next, nil
}
func (r *PostgresRepository) RevokeSession(ctx context.Context, tokenHash string) error {
	_, err := r.pool.Exec(ctx, `UPDATE sessions SET revoked_at=COALESCE(revoked_at,CURRENT_TIMESTAMP) WHERE token_hash=$1`, tokenHash)
	return err
}
func (r *PostgresRepository) SessionActive(ctx context.Context, sessionID string, now time.Time) (bool, error) {
	var active bool
	err := r.pool.QueryRow(ctx, `SELECT revoked_at IS NULL AND expires_at>$2 FROM sessions WHERE id=$1`, sessionID, now).Scan(&active)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("query active session: %w", err)
	}
	return active, nil
}
func (r *PostgresRepository) StorePasswordReset(ctx context.Context, userID, tokenHash string, expires time.Time) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO password_reset_tokens(user_id,token_hash,expires_at) VALUES($1,$2,$3)`, userID, tokenHash, expires)
	return err
}
func (r *PostgresRepository) ConsumeEmailVerification(ctx context.Context, tokenHash string, now time.Time) (string, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var userID string
	err = tx.QueryRow(ctx, `UPDATE email_verification_tokens SET consumed_at=$2 WHERE token_hash=$1 AND consumed_at IS NULL AND expires_at>$2 AND created_at + INTERVAL '120 seconds'>$2 RETURNING user_id`, tokenHash, now).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", ErrInvalidToken
	}
	if err != nil {
		return "", err
	}
	if _, err = tx.Exec(ctx, `UPDATE users SET email_verified_at=COALESCE(email_verified_at,$2),updated_at=$2 WHERE id=$1`, userID, now); err != nil {
		return "", err
	}
	if err = tx.Commit(ctx); err != nil {
		return "", err
	}
	return userID, nil
}
func (r *PostgresRepository) UpdatePassword(ctx context.Context, userID, passwordHash string, now time.Time) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err := tx.Exec(ctx, `UPDATE users SET password_hash=$2,updated_at=$3 WHERE id=$1`, userID, passwordHash, now)
	if err != nil {
		return err
	}
	if result.RowsAffected() != 1 {
		return ErrInvalidCredentials
	}
	if _, err = tx.Exec(ctx, `UPDATE sessions SET revoked_at=COALESCE(revoked_at,$2) WHERE user_id=$1`, userID, now); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (r *PostgresRepository) ResetPassword(ctx context.Context, tokenHash, passwordHash string, now time.Time) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	var userID string
	err = tx.QueryRow(ctx, `UPDATE password_reset_tokens SET consumed_at=$2 WHERE token_hash=$1 AND consumed_at IS NULL AND expires_at>$2 RETURNING user_id`, tokenHash, now).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrInvalidToken
	}
	if err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `UPDATE users SET password_hash=$2,updated_at=$3 WHERE id=$1`, userID, passwordHash, now); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `UPDATE sessions SET revoked_at=COALESCE(revoked_at,$2) WHERE user_id=$1`, userID, now); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
