package repository

import (
	"api-for-people/internal/user/model"
	"api-for-people/internal/user/service"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"strings"
)

type Repository struct {
	db *pgx.Conn
}

func NewRepository(db *pgx.Conn) service.RepositoryPostgres {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, person *model.Person) error {
	query := `INSERT INTO persons (
                     name, 
                     surname, 
                     patronymic, 
                     age, 
                     gender, 
                     nationality
                     ) 
				VALUES 
				    ($1, $2, $3, $4, $5, $6)
				    `
	_, err := r.db.Exec(ctx, query, &person.Name,
		person.Surname,
		person.Patronymic,
		person.Age,
		person.Gender,
		person.Nationality,
	)
	if err != nil {
		slog.Error("Error during exec command", err)
	}
	return nil
}
func (r *Repository) Get(ctx context.Context, id int) (model.Person, error) {
	query := `SELECT 
    	name, 
    	surname, 
    	patronymic, 
    	age, 
    	gender, 
    	nationality 
		FROM 
    		persons
    	WHERE 
    	    id = $1
		and
    	    isDeleted = false
`
	var person model.Person
	slog.Debug("Getting person from DB", id)
	err := r.db.QueryRow(ctx, query, id).Scan(&person.Name, &person.Surname, &person.Patronymic, &person.Age, &person.Gender, &person.Nationality)
	if err != nil {
		slog.Error("Error during query row", err)
	}
	return person, nil
}
func (r *Repository) GetAll(ctx context.Context, params model.UserQueryParams) ([]model.Person, error) {
	query := `SELECT name, surname, patronymic, age, gender, nationality FROM persons WHERE isDeleted=false`
	slog.Info("Getting all persons with params", params)
	var args []interface{}
	var whereClauses []string

	// Фильтрация
	if params.Name != "" {
		args = append(args, "%"+params.Name+"%")
		whereClauses = append(whereClauses, fmt.Sprintf("name LIKE $%d", len(args)))
	}
	if params.Surname != "" {
		args = append(args, "%"+params.Surname+"%")
		whereClauses = append(whereClauses, fmt.Sprintf("surname LIKE $%d", len(args)))
	}
	if params.Patronymic != "" {
		args = append(args, "%"+params.Patronymic+"%")
		whereClauses = append(whereClauses, fmt.Sprintf("patronymic LIKE $%d", len(args)))
	}
	if params.Gender != "" {
		args = append(args, params.Gender)
		whereClauses = append(whereClauses, fmt.Sprintf("gender = $%d", len(args)))
	}
	if params.Nationality != "" {
		args = append(args, params.Nationality)
		whereClauses = append(whereClauses, fmt.Sprintf("nationality = $%d", len(args)))
	}
	if params.Age > 0 {
		args = append(args, params.Age)
		whereClauses = append(whereClauses, fmt.Sprintf("age = $%d", len(args)))
	}

	// Правильное объединение условий
	if len(whereClauses) > 0 {
		query += " AND " + strings.Join(whereClauses, " AND ")
	}

	// Сортировка
	if params.SortBy != "" {
		sortOrder := "ASC"
		if params.SortOrder == "desc" {
			sortOrder = "DESC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", params.SortBy, sortOrder)
	}

	// Пагинация (только если указан Limit)
	if params.Limit > 0 {
		args = append(args, params.Limit, (params.Page-1)*params.Limit)
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))
	}
	slog.Info("Getting all persons", query, args)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		slog.Error("Error during query row", err)
	}
	defer rows.Close()

	persons := make([]model.Person, 0)
	for rows.Next() {
		var person model.Person
		if err = rows.Scan(
			&person.Name,
			&person.Surname,
			&person.Patronymic,
			&person.Age,
			&person.Gender,
			&person.Nationality,
		); err != nil {
			return nil, fmt.Errorf("failed to scan person: %w", err)
		}
		persons = append(persons, person)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return persons, nil
}
func (r *Repository) Update(ctx context.Context, person model.Person, id int) error {
	query := `UPDATE persons SET name = $1, surname = $2, patronymic = $3, age = $4, gender = $5, nationality = $6 WHERE id = $7 and isDeleted = false`
	slog.Debug("Updating person", id)
	tag, err := r.db.Exec(ctx, query, person.Name, person.Surname, person.Patronymic, person.Age, person.Gender, person.Nationality, id)
	if err != nil {
		slog.Error("Error during exec command", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("failed to update person: no rows affected")
	}
	return nil
}
func (r *Repository) Delete(ctx context.Context, id int) error {
	query := `UPDATE persons SET isDeleted = true WHERE id = $1`
	slog.Debug("Deleting person", id)
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		slog.Error("Error during exec command", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("failed to delete person: no rows affected")
	}
	return nil
}
