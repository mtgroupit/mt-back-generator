// +build integration

package migration

import (
	"context"
	"errors"
	"testing"

	"github.com/powerman/check"
	"github.com/pressly/goose"
)

var ErrDownNotSupported = errors.New("downgrade is not supported, restore from backup instead")

func TestDown(tt *testing.T) {
	t := check.T(tt)
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()
	db, err := connect(ctx, dbCfg)
	t.Must(t.Nil(err))
	defer db.Close()

	t.Must(t.Nil(Run(ctx, ".", "up", dbCfg)))
	for v, _ := goose.GetDBVersion(db); v > 0; v, _ = goose.GetDBVersion(db) {
		err := Run(ctx, ".", "down", dbCfg)
		if err != nil && t.Contains(err.Error(), ErrDownNotSupported.Error()) {
			t.Logf("downgrade from version %d is not supported", v)
			t.Nil(Run(ctx, ".", "up", dbCfg))
			return
		}
		t.Must(t.Nil(err))
		v2, err := goose.GetDBVersion(db)
		t.Nil(err)
		t.Less(v2, v)
	}
	v, err := goose.GetDBVersion(db)
	t.Nil(err)
	t.Zero(v)
	t.Nil(Run(ctx, ".", "up", dbCfg))
}
