package repository

import (
	"database/sql"
	"fmt"

	"github.com/KarmaBeLike/time-tracker-api/internal/models"
)

type UserRepository interface {
	GetUsers(page, limit int, filters map[string]string) ([]models.User, int, error)
	CreateUser(user *models.User) error
	DeleteUser(userId string) error
	UpdateUser(user *models.User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (name, surname, patronymic, passport_number, address)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at;
	`
	err := r.db.QueryRow(query, user.Name, user.Surname, user.Patronymic, user.PassportNumber, user.Address).Scan(&user.ID, &user.CreatedAt)
	return err
}

func (r *userRepository) GetUsers(page, limit int, filters map[string]string) ([]models.User, int, error) {
	var users []models.User
	query := "SELECT id, name, surname, patronymic, passport_number, task_ids,created_at FROM users WHERE 1=1"

	// Применение фильтров
	args := []interface{}{}
	name := ""
	query += " AND (name LIKE $1 OR $1 = '')"
	if nameq, ok := filters["name"]; ok {
		nameq = "%" + nameq + "%"
		name = nameq
	}
	args = append(args, name)

	passportNumber := ""
	query += " AND (passport_number = $2 OR $2 = '')"
	if passportNumberQ, ok := filters["passportNumber"]; ok {
		passportNumber = passportNumberQ
	}
	args = append(args, passportNumber)

	// Получение общего количества
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count_table", query)
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Применение пагинации и получение данных
	offset := (page - 1) * limit
	query += (" LIMIT $3 OFFSET $4")
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...) // Select fsf, fesfes From fsef where $3
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.PassportNumber, &user.TaskIDs, &user.CreatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) DeleteUser(userId string) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.Exec(query, userId)
	return err
}

func (r *userRepository) UpdateUser(user *models.User) error {
	query := `
		UPDATE users 
		SET name = $1, surname = $2, patronymic = $3, passport_number = $4, address = $5
		WHERE id = $6
	`
	_, err := r.db.Exec(query, user.Name, user.Surname, user.Patronymic, user.PassportNumber, user.Address, user.ID)
	return err
}
