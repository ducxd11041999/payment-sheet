package repository

import (
	"database/sql"
	"log"
)

type LogRepository struct {
	DB *sql.DB
}

func NewLogRepository(db *sql.DB) *LogRepository {
	return &LogRepository{DB: db}
}

func (r *LogRepository) Write(logEntry UserLog) error {
	_, err := r.DB.Exec(`
		INSERT INTO user_logs (username, method, path, ip_address, user_agent, body)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, logEntry.Username, logEntry.Method, logEntry.Path, logEntry.IPAddress, logEntry.UserAgent, logEntry.Body)

	if err != nil {
		log.Printf("Write user log: %v", err)
		return err
	}

	return nil
}

func (r *LogRepository) GetAllLogs() ([]UserLog, error) {
	rows, err := r.DB.Query(`SELECT id, username, method, path, ip_address, user_agent, body, 
       created_at FROM user_logs ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []UserLog
	for rows.Next() {
		var log UserLog
		err := rows.Scan(
			&log.ID,
			&log.Username,
			&log.Method,
			&log.Path,
			&log.IPAddress,
			&log.UserAgent,
			&log.Body,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}
