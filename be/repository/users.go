package repository

import "database/sql"

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetByUsername(username string) (*User, error) {
	row := r.DB.QueryRow(`SELECT id, username, password FROM users WHERE username = $1`, username)
	var u User
	err := row.Scan(&u.ID, &u.Username, &u.Password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Create(user *User) error {
	_, err := r.DB.Exec(`INSERT INTO users (id, username, password) VALUES ($1, $2, $3)`,
		user.ID, user.Username, user.Password)
	return err
}
