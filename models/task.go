package models

import "database/sql"

// ProjectTask object
type ProjectTask struct {
	ID             int    `json:"id"`
	UserID         int    `json:"userId"`
	ProjectID      int    `json:"projectId"`
	Title          string `json:"title"`
	IsCompleted    int    `json:"isCompleted"`
	UserName       string `json:"userName"`
	ProjectName    string `json:"projectName"`
	AttachmentPath string `json:"attachmentPath"`
}

// GetTasks is used to retrieve all project task
func GetTasks(db *sql.DB) ([]ProjectTask, error) {
	rows, err := db.Query("SELECT pt.id, pt.title, pt.is_completed, u.username, p.name, " +
		"pt.attachment_path FROM project_task AS pt INNER JOIN users AS u ON u.id = pt.user_id " +
		"INNER JOIN projects AS p ON p.id = pt.project_id;")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	ProjectTasks := []ProjectTask{}

	for rows.Next() {
		var t ProjectTask
		if err := rows.Scan(&t.ID, &t.Title, &t.IsCompleted, &t.UserName, &t.ProjectName, &t.AttachmentPath); err != nil {
			return nil, err
		}
		ProjectTasks = append(ProjectTasks, t)
	}

	return ProjectTasks, nil
}

// GetCompletedTasks will retrieve all the completed task
func GetCompletedTasks(db *sql.DB) ([]ProjectTask, error) {
	rows, err := db.Query("SELECT pt.id, pt.title, pt.is_completed, u.username, p.name, " +
		"pt.attachment_path FROM project_task AS pt INNER JOIN users AS u ON u.id = pt.user_id " +
		"INNER JOIN projects AS p ON p.id = pt.project_id " +
		"WHERE is_completed = 1;")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	ProjectTasks := []ProjectTask{}

	for rows.Next() {
		var t ProjectTask
		if err := rows.Scan(&t.ID, &t.Title, &t.IsCompleted, &t.UserName, &t.ProjectName, &t.AttachmentPath); err != nil {
			return nil, err
		}
		ProjectTasks = append(ProjectTasks, t)
	}

	return ProjectTasks, nil
}

// GetTaskByID is used to retrieve a single project task
func (p *ProjectTask) GetTaskByID(db *sql.DB) error {
	return db.QueryRow("SELECT id,user_id,project_id,title,is_completed,attachment_path FROM project_task WHERE id=$1",
		p.ID).Scan(&p.ID, &p.UserID, &p.ProjectID, &p.Title, &p.IsCompleted, &p.AttachmentPath)
}

// GetTasksByUserID is used to retrieve all project task for a particular user and status
func (p *ProjectTask) GetTasksByUserID(db *sql.DB, taskStatus int) ([]ProjectTask, error) {
	rows, err := db.Query("SELECT pt.id, pt.title, pt.is_completed, u.username, u.id, p.name, "+
		"pt.attachment_path FROM project_task AS pt INNER JOIN users AS u ON u.id = pt.user_id "+
		"INNER JOIN projects AS p ON p.id = pt.project_id "+
		"WHERE pt.user_id = $1 and is_completed = $2", p.UserID, taskStatus)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	ProjectTasks := []ProjectTask{}

	for rows.Next() {
		var t ProjectTask
		if err := rows.Scan(&t.ID, &t.Title, &t.IsCompleted, &t.UserName, &t.UserID, &t.ProjectName, &t.AttachmentPath); err != nil {
			return nil, err
		}
		ProjectTasks = append(ProjectTasks, t)
	}

	return ProjectTasks, nil
}

// GetTasksForUser is used to retrieve all project task for a particular user open and completed
func (p *ProjectTask) GetTasksForUser(db *sql.DB) ([]ProjectTask, error) {
	rows, err := db.Query("SELECT pt.id, pt.title, pt.is_completed, u.username, u.id, p.name, "+
		"pt.attachment_path FROM project_task AS pt INNER JOIN users AS u ON u.id = pt.user_id "+
		"INNER JOIN projects AS p ON p.id = pt.project_id "+
		"WHERE pt.user_id = $1", p.UserID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	ProjectTasks := []ProjectTask{}

	for rows.Next() {
		var t ProjectTask
		if err := rows.Scan(&t.ID, &t.Title, &t.IsCompleted, &t.UserName, &t.UserID, &t.ProjectName, &t.AttachmentPath); err != nil {
			return nil, err
		}
		ProjectTasks = append(ProjectTasks, t)
	}

	return ProjectTasks, nil
}

// UpdateTask updates and existing project task
func (p *ProjectTask) UpdateTask(db *sql.DB) error {
	_, err := db.Exec("Update project_task SET user_id=$1,project_id=$2,title=$3,is_completed=$4,attachment_path=$5 WHERE id=$6", p.UserID, p.ProjectID, p.Title, p.IsCompleted, p.AttachmentPath, p.ID)
	return err
}

// CreateTask creates a new project task
func (p *ProjectTask) CreateTask(db *sql.DB) error {
	_, err := db.Exec("INSERT INTO project_task (user_id,project_id,title,is_completed,attachment_path) VALUES ($1,$2,$3,$4,$5)", p.UserID, p.ProjectID, p.Title, p.IsCompleted, p.AttachmentPath)
	return err
}

// DeleteTaskByID is used to delete a single project task
func (p *ProjectTask) DeleteTaskByID(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM project_task WHERE id=$1", p.ID)
	return err
}
