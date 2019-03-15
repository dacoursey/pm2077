package models

import "database/sql"

// Contact object
type Contact struct {
	ID     int    `json:"id"`
	CustID int    `json:"cust_id"`
	FName  string `json:"fname"`
	LName  string `json:"lname"`
	Email  string `json:"email"`
}

// GetContact is used to retrieve one contact object from the db.
func (p *Contact) GetContact(db *sql.DB) error {
	return db.QueryRow("SELECT cust_id, fname, lname, email FROM contacts WHERE id=$1",
		p.ID).Scan(&p.CustID, &p.FName, &p.LName, &p.Email)
}

// UpdateContact is used to write changes to one contact object in the db.
func (p *Contact) UpdateContact(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE contacts SET cust_id=$1, fname=$2, lname=$3, email=$4 WHERE id=$5",
			p.CustID, p.FName, p.LName, p.Email, p.ID)

	return err
}

// DeleteContact is used to remove one contact object from the db.
func (p *Contact) DeleteContact(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM contacts WHERE id=$1", p.ID)

	return err
}

// CreateContact is used to add one new contact object to the db.
func (p *Contact) CreateContact(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO contacts(cust_id, fname, lname, email) VALUES($1, $2, $3, $4) RETURNING id",
		p.CustID, p.FName, p.LName, p.Email).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}

// GetContacts is used to retrieve all contact objects from the db.
func GetContacts(db *sql.DB, start, count int) ([]Contact, error) {
	rows, err := db.Query(
		"SELECT * FROM contacts LIMIT $1 OFFSET $2", count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	Contacts := []Contact{}

	for rows.Next() {
		var p Contact
		if err := rows.Scan(&p.ID, &p.CustID, &p.FName, &p.LName, &p.Email); err != nil {
			return nil, err
		}
		Contacts = append(Contacts, p)
	}

	return Contacts, nil
}
