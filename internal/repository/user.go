package repositories

import (
	"database/sql"
	"fmt"

	"github.com/KarmaBeLike/time-tracker-api/internal/models"
)

type UserRepository interface {
	GetUsers(page, limit int, filters map[string]string) ([]models.User, int, error)
}

type userRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) GetUsers(page, limit int, filters map[string]string) ([]models.User, int, error) {
	var users []models.User
	query := "SELECT id, name, passport_number FROM users WHERE 1=1"

	// Применение фильтров
	args := []interface{}{}
	if name, ok := filters["name"]; ok {
		query += " AND name LIKE $1"
		args = append(args, "%"+name+"%")
	}
	if passportNumber, ok := filters["passportNumber"]; ok {
		query += " AND passport_number = $2"
		args = append(args, passportNumber)
	}

	// Получение общего количества
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count_table", query)
	err := r.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Применение пагинации и получение данных
	offset := (page - 1) * limit
	query += (" LIMIT $3 OFFSET $4")
	args = append(args, limit, offset)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.PassportNumber); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
