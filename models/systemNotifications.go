package models

import "database/sql"

// SystemNotification object
type SystemNotification struct {
	ID      int32    `json:"id"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

// GetNotifications is used to retrieve all system notifications
func GetNotifications(db *sql.DB) ([]SystemNotification, error) {
	rows, err := db.Query("SELECT * FROM system_notifications")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	SystemNotifications := []SystemNotification{}

	for rows.Next() {
		var n SystemNotification
		if err := rows.Scan(&n.ID, &n.Title, &n.Message); err != nil {
			return nil, err
		}
		SystemNotifications = append(SystemNotifications, n)
	}

	return SystemNotifications, nil
}

// GetNotificationByID is used to retrieve a single system notification
func (n *SystemNotification) GetNotificationByID(db *sql.DB) error {
	return db.QueryRow("SELECT id, title, message FROM system_notifications WHERE id=$1",
		n.ID).Scan(&n.ID, &n.Title, &n.Message)
}

// UpdateNotificationByID updates a system notification by the ID
func (n *SystemNotification) UpdateNotificationByID(db *sql.DB) error {
	_, err := db.Exec("UPDATE system_notifications SET title=$1, message=$2 WHERE id=$3", n.Title, n.Message, n.ID)

	return err
}

// CreateNotification creates a new system notification
func (n *SystemNotification) CreateNotification(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO system_notifications (title,message) VALUES ($1,$2)", n.Title, n.Message)

	return err
}

// DeleteNotificationByID is used to delete a single system notification
func (n *SystemNotification) DeleteNotificationByID(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM system_notifications WHERE id=$1", n.ID)
	return err
}