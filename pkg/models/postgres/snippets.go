package postgres

import (
	"database/sql"
	"github.com/nikishin42/toDoList/pkg/models"
	"time"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	result, err := m.DB.Exec(`insert into snippets (title, content, created, expires) values ($1, $2, now, interval $3 days)`,
		title, content, time.Now(), time.Now().AddDate(0, 0, expires))
	if err != nil {
		return 0, nil
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return nil, nil
}
