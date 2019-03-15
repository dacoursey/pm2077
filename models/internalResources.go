package models

import "database/sql"

// Resource object
type Resource struct {
	ID      int32  `json:"id"`
	Title   string `json:"title"`
	URL		string `json:"url"`
}

// GetResources is used to retrieve all resource
func GetResources(db *sql.DB) ([]Resource, error) {
	rows, err := db.Query("SELECT * FROM internal_resources")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	Resources := []Resource{}

	for rows.Next() {
		var n Resource
		if err := rows.Scan(&n.ID, &n.Title, &n.URL); err != nil {
			return nil, err
		}
		Resources = append(Resources, n)
	}

	return Resources, nil
}

// GetResourceByID is used to retrieve a single resource
func (n *Resource) GetResourceByID(db *sql.DB) error {
	return db.QueryRow("SELECT id, title, url FROM internal_resources WHERE id=$1",
		n.ID).Scan(&n.ID, &n.Title, &n.URL)
}

// UpdateResourceByID updates a resource by the ID
func (n *Resource) UpdateResourceByID(db *sql.DB) error {
	_, err := db.Exec("UPDATE internal_resources SET title=$1, URL=$2 WHERE id=$3", n.Title, n.URL, n.ID)

	return err
}

// CreateResource creates a new resource
func (n *Resource) CreateResource(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO internal_resources (title,url) VALUES ($1,$2)", n.Title, n.URL)

	return err
}

// DeleteResourceByID is used to delete a single resource
func (n *Resource) DeleteResourceByID(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM internal_resources WHERE id=$1", n.ID)
	return err
}