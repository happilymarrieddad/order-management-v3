package repos

import (
	"context"

	"github.com/happilymarrieddad/order-management-v3/api/types"
	"golang.org/x/crypto/bcrypt"
	"xorm.io/xorm"
)

// UsersRepo defines the interface for user data operations.
//
//go:generate mockgen -source=./users.go -destination=./mocks/users.go -package=mock_repos UsersRepo
type UsersRepo interface {
	Find(ctx context.Context, opts *UserFindOpts) ([]*types.User, int64, error)
	Get(ctx context.Context, id int64) (*types.User, bool, error)
	GetByEmail(ctx context.Context, email string) (*types.User, bool, error)
	Create(ctx context.Context, user *types.User) error
	CreateTx(ctx context.Context, tx *xorm.Session, user *types.User) error
	UpdateTx(ctx context.Context, tx *xorm.Session, user *types.User) error
	Update(ctx context.Context, user *types.User) error
	DeleteTx(ctx context.Context, tx *xorm.Session, id int64) error
	Delete(ctx context.Context, id int64) error
	UpdatePassword(ctx context.Context, userID int64, newPassword string) error
	UpdatePasswordTx(ctx context.Context, tx *xorm.Session, userID int64, newPassword string) error
	GetIncludeInvisible(ctx context.Context, id int64) (*types.User, bool, error)
}

type usersRepo struct {
	db *xorm.Engine
}

// NewUsersRepo creates a new UsersRepo.
func NewUsersRepo(db *xorm.Engine) UsersRepo {
	return &usersRepo{db: db}
}

// UserFindOpts provides options for finding users.
type UserFindOpts struct {
	IDs    []int64
	Emails []string
	Limit  int
	Offset int
}

// Create inserts a new user into the database.
func (r *usersRepo) Create(ctx context.Context, user *types.User) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.CreateTx(ctx, tx, user)
	})
	return err
}

func (r *usersRepo) CreateTx(ctx context.Context, tx *xorm.Session, user *types.User) error {
	if err := types.Validate(user); err != nil {
		return err
	}

	// Hash the user's password before saving to the database.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.Visible = true

	_, err = tx.Context(ctx).Insert(user)
	return err
}

// Delete performs a soft delete on a user.
func (r *usersRepo) Delete(ctx context.Context, id int64) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.DeleteTx(ctx, tx, id)
	})
	return err
}

// DeleteTx performs a soft delete on a user by setting their visible flag to false.
func (r *usersRepo) DeleteTx(ctx context.Context, tx *xorm.Session, id int64) error {
	_, err := tx.Context(ctx).ID(id).Cols("visible").Update(&types.User{Visible: false})
	return err
}

// Find retrieves a list of visible users with pagination and filtering, and a total count.
func (r *usersRepo) Find(ctx context.Context, opts *UserFindOpts) ([]*types.User, int64, error) {
	s := r.db.NewSession().Context(ctx)
	defer s.Close()
	s.Where("visible = ?", true)
	applyUserFindOpts(s, opts)
	var users []*types.User
	count, err := s.FindAndCount(&users)
	return users, count, err
}

// Get retrieves a single visible user by their ID.
func (r *usersRepo) Get(ctx context.Context, id int64) (*types.User, bool, error) {
	user := new(types.User)
	has, err := r.db.Context(ctx).ID(id).Where("visible = ?", true).Get(user)
	return user, has, err
}

// GetIncludeInvisible retrieves a single user by their ID, regardless of visibility.
func (r *usersRepo) GetIncludeInvisible(ctx context.Context, id int64) (*types.User, bool, error) {
	user := new(types.User)
	has, err := r.db.Context(ctx).ID(id).Get(user)
	return user, has, err
}

// GetByEmail retrieves a single user by their email.
func (r *usersRepo) GetByEmail(ctx context.Context, email string) (*types.User, bool, error) {
	user := new(types.User)
	has, err := r.db.Context(ctx).Where("email = ?", email).Get(user)
	return user, has, err
}

// Update modifies an existing user in the database.
func (r *usersRepo) Update(ctx context.Context, user *types.User) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.UpdateTx(ctx, tx, user)
	})
	return err
}

// UpdateTx modifies an existing user's non-sensitive details within a transaction.
// This function explicitly does NOT update the password.
func (r *usersRepo) UpdateTx(ctx context.Context, tx *xorm.Session, user *types.User) error {
	if err := types.Validate(user); err != nil {
		return err
	}

	// Explicitly update only non-sensitive fields. Password must be updated via UpdatePassword.
	_, err := tx.Context(ctx).ID(user.ID).Cols("first_name", "last_name", "address_id", "roles").Update(user)
	return err
}

// UpdatePassword updates a user's password.
func (r *usersRepo) UpdatePassword(ctx context.Context, userID int64, newPassword string) error {
	_, err := wrapInSession(r.db, func(tx *xorm.Session) (*struct{}, error) {
		return nil, r.UpdatePasswordTx(ctx, tx, userID, newPassword)
	})
	return err
}

// UpdatePasswordTx updates a user's password within a transaction.
func (r *usersRepo) UpdatePasswordTx(ctx context.Context, tx *xorm.Session, userID int64, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update only the password column for the given user ID.
	_, err = tx.Context(ctx).ID(userID).Cols("password").Update(&types.User{Password: string(hashedPassword)})
	return err
}

// applyUserFindOpts is a helper function to build the query based on find options.
func applyUserFindOpts(s *xorm.Session, opts *UserFindOpts) {
	if opts == nil {
		return
	}
	if len(opts.IDs) > 0 {
		s.In("id", opts.IDs)
	}
	if len(opts.Emails) > 0 {
		s.In("email", opts.Emails)
	}
	if opts.Limit > 0 {
		s.Limit(opts.Limit, opts.Offset)
	}
}
