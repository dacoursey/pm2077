package models

import "database/sql"

// Role Applicaiton Roles
type Role struct {
	ID		int `json:"id"`
	Name	string `json:"name"`
}

// GetAppRoles is used to retrieve all application roles
func GetAppRoles(db *sql.DB) ([]Role, error) {
	rows, err := db.Query("SELECT id, name FROM roles")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	Roles := []Role{}

	for rows.Next() {
		var r Role
		if err := rows.Scan(&r.ID, &r.Name); err != nil {
			return nil, err
		}
		Roles = append(Roles, r)
	}

	return Roles, nil
}