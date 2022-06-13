package command

import (
	"github.com/pkg/errors"
	"github.com/rlawnsxo131/madre-server-v3/external/datastore/rdb"
	"github.com/rlawnsxo131/madre-server-v3/internal/domain/account"
)

type socialAccountCommandRepository struct {
	db rdb.Database
}

func NewSocialAccountCommandRepository(db rdb.Database) account.SocialAccountCommandRepository {
	return &socialAccountCommandRepository{db}
}

func (r *socialAccountCommandRepository) Create(sa *account.SocialAccount) (string, error) {
	var id string

	query := "INSERT INTO social_account(user_id, provider, social_id)" +
		" VALUES(:user_id, :provider, :social_id)" +
		" RETURNING id"

	err := r.db.PrepareNamedGet(
		&id,
		query,
		sa,
	)
	if err != nil {
		return "", errors.Wrap(err, "socialaccount WriteRepository create")
	}

	return id, err
}
