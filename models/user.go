package models

import "database/sql"

// User object
type User struct {
	ID       int    `json:"id"`
	Username string `json:"name"`
	Password string `json:"password"`
	Role	 string `json:"role"`
}

// GetUser is used to retrieve one user object from the db.
func (p *User) GetUser(db *sql.DB) error {
	return db.QueryRow("SELECT username, password, role FROM users WHERE id=$1",p.ID).Scan(&p.Username, &p.Password, &p.Role)
}

// GetUserByCreds is used to retrieve one user object from the db for authentication.
func (p *User) GetUserByCreds(db *sql.DB) error {
	return db.QueryRow("SELECT u.id, u.username, u.password, r.name FROM users as u INNER JOIN roles as r ON r.id = u.role WHERE username=$1",
		p.Username).Scan(&p.ID, &p.Username, &p.Password, &p.Role)
}

// UpdateUser is used to write changes to one user object in the db.
func (p *User) UpdateUser(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE users SET username=$1, password=$2, role=$3 where id=$4", p.Username, p.Password, p.Role, p.ID)

	return err
}

// DeleteUser is used to remove one user object from the db.
func (p *User) DeleteUser(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM users WHERE id=$1", p.ID)

	return err
}

// CreateUser is used to add one new user object to the db.
func (p *User) CreateUser(db *sql.DB) error {
	err := db.QueryRow("INSERT INTO users(username, password, role) VALUES($1, $2, $3) RETURNING id", p.Username, p.Password, p.Role).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}

// GetUsers is used to retrieve all user objects from the db.
func GetUsers(db *sql.DB, start, count int) ([]User, error) {
	rows, err := db.Query("SELECT u.id, u.username, u.password, r.name " +
	                      "FROM users as u INNER JOIN roles as r ON r.id = u.role LIMIT $1 OFFSET $2", count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	Users := []User{}

	for rows.Next() {
		var p User
		if err := rows.Scan(&p.ID, &p.Username, &p.Password, &p.Role); err != nil {
			return nil, err
		}
		Users = append(Users, p)
	}

	return Users, nil
}