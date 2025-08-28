package repos

import (
	"log"

	"github.com/happilymarrieddad/order-management-v3/api/types"
	"xorm.io/xorm"
)

func handleRollback(sesh *xorm.Session, err error) error {
	if rollBackErr := sesh.Rollback(); rollBackErr != nil {
		log.Printf("unable to rollback with err: %s", rollBackErr.Error())
	}

	return err
}

func wrapInSession[T any](db *xorm.Engine, fn func(*xorm.Session) (*T, error)) (*T, error) {
	session := db.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return nil, types.NewInternalServerError("unable to start transaction with err: " + err.Error())
	}

	res, err := fn(session)
	if err != nil {
		return nil, handleRollback(session, err)
	}

	if err := session.Commit(); err != nil {
		return nil, handleRollback(session, err)
	}

	return res, nil
}
