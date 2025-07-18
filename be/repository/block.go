package repository

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type BlockRepository struct {
	DB *sql.DB
}

func NewBlockRepository(db *sql.DB) *BlockRepository {
	return &BlockRepository{
		DB: db,
	}
}

func (r *BlockRepository) GetIDByMonth(month string) (string, bool, error) {
	row := r.DB.QueryRow(`SELECT id, locked FROM blocks WHERE month = $1`, month)
	var blockID string
	var locked bool
	if err := row.Scan(&blockID, &locked); err != nil {
		return "", false, fiber.ErrNotFound
	}

	return blockID, locked, nil
}

func (r *BlockRepository) Get(id string) (string, bool, error) {
	row := r.DB.QueryRow(`SELECT id, locked FROM blocks WHERE id = $1`, id)
	var blockID string
	var locked bool
	if err := row.Scan(&blockID, &locked); err != nil {
		return "", false, fiber.ErrNotFound
	}

	return blockID, locked, nil
}

func (r *BlockRepository) Lock(month string) error {
	_, err := r.DB.Exec(`UPDATE blocks SET locked = true WHERE month = $1`, month)
	return err
}

func (r *BlockRepository) Unlock(month string) error {
	_, err := r.DB.Exec(`UPDATE blocks SET locked = false WHERE month = $1`, month)
	return err
}

func (r *BlockRepository) Create(block Block) error {
	_, err := r.DB.Exec(`INSERT INTO blocks (id, month, locked) VALUES ($1, $2, $3)`, block.ID, block.Month, block.Locked)
	if err != nil {
		return err
	}
	stmt, err := r.DB.Prepare(`INSERT INTO members (id, block_id, name, ratio, debt) VALUES ($1, $2, $3, $4, $5)`)
	if err != nil {
		return err
	}
	for _, m := range block.Members {
		m.ID = uuid.New().String()
		if _, err := stmt.Exec(m.ID, block.ID, m.Name, m.Ratio, m.Debt); err != nil {
			return err
		}
	}
	return nil
}

func (r *BlockRepository) GetAllBlocks() ([]Block, error) {
	rows, err := r.DB.Query(`SELECT id, month, locked FROM blocks ORDER BY month DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []Block
	for rows.Next() {
		var b Block
		if err := rows.Scan(&b.ID, &b.Month, &b.Locked); err != nil {
			return nil, err
		}
		blocks = append(blocks, b)
	}
	return blocks, nil
}

func (r *BlockRepository) DeleteBlock(blockID string) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	// Xóa chi tiết giao dịch (transaction_details)
	_, err = tx.Exec(`
        DELETE FROM transaction_details 
        WHERE transaction_id IN (
            SELECT id FROM transactions WHERE block_id = $1
        )`, blockID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Xoá transactions liên quan
	_, err = tx.Exec("DELETE FROM transactions WHERE block_id = $1", blockID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Xoá members liên quan
	_, err = tx.Exec("DELETE FROM members WHERE block_id = $1", blockID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Xoá block
	_, err = tx.Exec("DELETE FROM blocks WHERE id = $1", blockID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
