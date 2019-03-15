package models

import "database/sql"

// Customer object
type Customer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	POC  int    `json:"poc"`
}

// GetCustomer is used to retrieve one customer object from the db.
func (p *Customer) GetCustomer(db *sql.DB) error {
	return db.QueryRow("SELECT name, poc FROM customers WHERE id=$1",
		p.ID).Scan(&p.Name, &p.POC)
}

// UpdateCustomer is used to write changes to one customer object in the db.
func (p *Customer) UpdateCustomer(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE customers SET name=$1, poc=$2 WHERE id=$3", p.Name, p.POC, p.ID)

	return err
}

// DeleteCustomer is used to remove one customer object from the db.
func (p *Customer) DeleteCustomer(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM customers WHERE id=$1", p.ID)

	return err
}

// CreateCustomer is used to add one new customer object to the db.
func (p *Customer) CreateCustomer(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO customers(name, poc) VALUES($1, $2) RETURNING id", p.Name, p.POC).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}

// GetCustomers is used to retrieve all customer objects from the db.
func GetCustomers(db *sql.DB, start, count int) ([]Customer, error) {
	rows, err := db.Query(
		"SELECT * FROM customers LIMIT $1 OFFSET $2", count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	Customers := []Customer{}

	for rows.Next() {
		var p Customer
		if err := rows.Scan(&p.ID, &p.Name, &p.POC); err != nil {
			return nil, err
		}
		Customers = append(Customers, p)
	}

	return Customers, nil
}
