package models

import "database/sql"

// Project object
type Project struct {
	ID        int    `json:"id"`
	CustID    int    `json:"cust_id"`
	Name      string `json:"name"`
	StartDate string `json:"start_date"`
	Hours     int    `json:"hours"`
}

// GetProject is used to retrieve one customer object from the db.
func (p *Project) GetProject(db *sql.DB) error {
	return db.QueryRow("SELECT cust_id, name, start_date, hours FROM projects WHERE id=$1",
		p.ID).Scan(&p.CustID, &p.Name, &p.StartDate, &p.Hours)
}

// UpdateProject is used to write changes to one customer object in the db.
func (p *Project) UpdateProject(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE projects SET cust_id=$1, name=$2, start_date=$3, hours=$4 where id=$5", p.CustID, p.Name, p.StartDate, p.Hours, p.ID)

	return err
}

// DeleteProject is used to remove one customer object from the db.
func (p *Project) DeleteProject(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM projects WHERE id=$1", p.ID)

	return err
}

// CreateProject is used to add one new customer object to the db.
func (p *Project) CreateProject(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO projects(cust_id, name, start_date, hours) VALUES($1, $2, $3, $4) RETURNING id",
		p.CustID, p.Name, p.StartDate, p.Hours).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}

// GetProjects is used to retrieve all customer objects from the db.
func GetProjects(db *sql.DB, start, count int) ([]Project, error) {
	rows, err := db.Query(
		"SELECT * FROM projects LIMIT $1 OFFSET $2", count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	Projects := []Project{}

	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.CustID, &p.Name, &p.StartDate, &p.Hours); err != nil {
			return nil, err
		}
		Projects = append(Projects, p)
	}

	return Projects, nil
}
