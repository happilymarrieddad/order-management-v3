package repos

import (
	"fmt"
	"log"

	"xorm.io/xorm"
)

func handleRollback(sesh *xorm.Session, err error) error {
	if rollBackErr := sesh.Rollback(); rollBackErr != nil {
		log.Printf("unable to rollback with err: %s", rollBackErr.Error())
	}

	return err
}

// wrapInSession is a helper function to manage database transactions.
// It starts a new session, begins a transaction, executes the provided function,
// and then commits or rolls back based on whether the function returns an error.
func wrapInSession[T any](db *xorm.Engine, fn func(tx *xorm.Session) (T, error)) (T, error) {
	var (
		res T
		err error
	)

	sess := db.NewSession()
	defer sess.Close()

	if err = sess.Begin(); err != nil {
		return res, fmt.Errorf("failed to begin transaction: %w", err)
	}

	res, err = fn(sess)
	if err != nil {
		sess.Rollback() // Attempt to rollback, but return the original error.
		return res, err
	}

	return res, sess.Commit()
}
