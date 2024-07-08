package repository

import (
	"database/sql"
	"time"

	"github.com/KarmaBeLike/time-tracker-api/internal/models"
)

type TaskRepository interface {
	GetWorklogs(userId int, startDate, endDate string) ([]models.Task, error)
	StartTask(userId, taskId int, startTime time.Time) error
	StopTask(userId, taskId int, endTime time.Time) error
}

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) GetWorklogs(userId int, startDate, endDate string) ([]models.Task, error) {
	query := `
		SELECT 
			t.id, t.name, t.description, SUM(EXTRACT(EPOCH FROM (w.end_time - w.start_time))/3600) as total_hours, 
			SUM(EXTRACT(EPOCH FROM (w.end_time - w.start_time))/60) % 60 as total_minutes
		FROM 
			tasks t
		JOIN 
			worklogs w ON t.id = w.task_id
		WHERE 
			w.user_id = $1 AND w.start_time >= $2 AND w.end_time <= $3
		GROUP BY 
			t.id
		ORDER BY 
			total_hours DESC, total_minutes DESC;
	`

	rows, err := r.db.Query(query, userId, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Name, &task.Description, &task.TotalHours, &task.TotalMinutes)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *taskRepository) StartTask(userId, taskId int, startTime time.Time) error {
	query := `
		INSERT INTO worklogs (user_id, task_id, start_time)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.Exec(query, userId, taskId, startTime)
	return err
}

func (r *taskRepository) StopTask(userId, taskId int, endTime time.Time) error {
	query := `
		UPDATE worklogs
		SET end_time = $1
		WHERE user_id = $2 AND task_id = $3 AND end_time IS NULL
	`
	_, err := r.db.Exec(query, endTime, userId, taskId)
	return err
}
