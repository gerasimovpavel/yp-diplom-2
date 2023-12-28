package repoPostgres

import (
	"context"
	"errors"
	trmpgx "github.com/avito-tech/go-transaction-manager/pgxv5"
	trmman "github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
	"yp-diplom-2/cmd/internal/domain"
)

type UserRepo struct {
	pool    *pgxpool.Pool
	getter  *trmpgx.CtxGetter
	manager *trmman.Manager
}

func NewUserRepo(pool *pgxpool.Pool, getter *trmpgx.CtxGetter) *UserRepo {
	return &UserRepo{
		pool:    pool,
		getter:  getter,
		manager: trmman.Must(trmpgx.NewDefaultFactory(pool)),
	}
}

func (r *UserRepo) Create(ctx context.Context, user domain.User) error {

	err := r.manager.Do(ctx, func(ctx context.Context) error {
		conn := r.getter.DefaultTrOrDB(ctx, r.pool)
		sqlShmt := `INSERT INTO users (user_id, last_name, first_name, email, password)
				VALUES (@user_id, @last_name, @first_name, @email, @password)`

		_, err := conn.Exec(ctx, sqlShmt, pgx.NamedArgs{
			"user_id":    user.UserID,
			"last_name":  user.LastName,
			"first_name": user.FirstName,
			"email":      user.Email,
			"password":   user.Password,
		})

		if err != nil {
			if IsDuplicate(err) {
				return domain.ErrUserAlreadyExists
			}
			return err
		}

		sqlShmt = `INSERT INTO verifications (user_id, code, verifed_at) 
			   VALUES (@user_id, @code, NULL)`
		_, err = conn.Exec(ctx, sqlShmt, pgx.NamedArgs{
			"user_id": user.UserID,
			"code": gofakeit.Password(
				true,
				true,
				true,
				false,
				false,
				32),
		})

		if err != nil {
			if IsDuplicate(err) {
				return domain.ErrUserAlreadyExists
			}
			return err
		}
		return err
	})

	return err
}

func (r *UserRepo) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	var user domain.User
	conn := r.getter.DefaultTrOrDB(ctx, r.pool)
	sqlShmt := `SELECT * FROM users WHERE email=@email AND password=@password`

	rows, err := conn.Query(ctx, sqlShmt, pgx.NamedArgs{
		"email":    email,
		"password": password,
	})
	if err != nil {
		return user, err
	}
	err = pgxscan.ScanOne(&user, rows)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, domain.ErrUserNotFound
		}
		return user, err
	}
	return user, err
}

func (r *UserRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error) {
	var user domain.User
	conn := r.getter.DefaultTrOrDB(ctx, r.pool)
	sqlShmt := `	select
						u.*
					from
						users u
					inner join sessions s on
						s.user_id = u.user_id
					where
						s.refresh_token = $1`
	err := pgxscan.Get(ctx, conn, &user, sqlShmt, refreshToken)
	if err != nil {
		return user, err
	}
	return user, err
}

func (r *UserRepo) Verify(ctx context.Context, userID uuid.UUID, code string) error {
	conn := r.getter.DefaultTrOrDB(ctx, r.pool)
	sqlShmt := `	UPDATE verifications
				    SET code='', verifed_at = @verifed_at
					WHERE user_id=@user_id AND code=@code`

	tag, err := conn.Exec(ctx, sqlShmt, pgx.NamedArgs{
		"verifed_at": time.Now(),
		"user_id":    userID,
		"code":       "",
	})
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrVerificationCodeInvalid
	}

	return nil
}

func (r *UserRepo) SetSession(ctx context.Context, userID uuid.UUID, session domain.Session) error {
	conn := r.getter.DefaultTrOrDB(ctx, r.pool)

	sqlShmt := `DELETE FROM sessions WHERE user_id=$1`
	conn.Exec(ctx, sqlShmt, userID)

	sqlShmt = `	INSERT INTO sessions (user_id, refresh_token, expires_at)
				VALUES (@user_id, @refresh_token, @expires_at)`
	_, err := conn.Exec(ctx, sqlShmt, pgx.NamedArgs{
		"user_id":       userID,
		"refresh_token": session.RefreshToken,
		"expires_at":    session.ExpiresAt,
	})
	return err
}
