package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type TransactionRepository struct {
	DB *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{DB: db}
}

func (r *TransactionRepository) GetByBlockID(blockID string) ([]Transaction, error) {
	rows, err := r.DB.Query(`SELECT id, description, amount, payer, created_at, ratios FROM transactions WHERE block_id = $1`, blockID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []Transaction
	for rows.Next() {
		var tx Transaction
		var ratiosJSON []byte

		tx.BlockID = blockID
		tx.Details = map[string]float64{}
		tx.Ratios = map[string]float64{}

		err := rows.Scan(&tx.ID, &tx.Description, &tx.Amount, &tx.Payer, &tx.CreatedAt, &ratiosJSON)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(ratiosJSON, &tx.Ratios); err != nil {
			return nil, err
		}

		// Query details
		detailRows, err := r.DB.Query(`SELECT member_id, amount FROM transaction_details WHERE transaction_id = $1`, tx.ID)
		if err != nil {
			return nil, err
		}
		for detailRows.Next() {
			var memberID string
			var amount float64
			if err := detailRows.Scan(&memberID, &amount); err != nil {
				detailRows.Close()
				return nil, err
			}
			tx.Details[memberID] = amount
		}
		detailRows.Close()

		txs = append(txs, tx)
	}
	return txs, nil
}

func (r *TransactionRepository) Add(tx Transaction) error {
	ratiosJSON, err := json.Marshal(tx.Ratios)
	if err != nil {
		return err
	}

	_, err = r.DB.Exec(`
		INSERT INTO transactions (id, block_id, payer, amount, description, created_at, ratios)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, tx.ID, tx.BlockID, tx.Payer, tx.Amount, tx.Description, tx.CreatedAt, ratiosJSON)

	return err
}

func (r *TransactionRepository) AddDetails(txID string, details map[string]float64) error {
	stmt, err := r.DB.Prepare(`
		INSERT INTO transaction_details (transaction_id, member_id, amount)
		VALUES ($1, $2, $3)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for memberID, amount := range details {
		if _, err := stmt.Exec(txID, memberID, amount); err != nil {
			return err
		}
	}
	return nil
}

func (r *TransactionRepository) GetByID(id string) (Transaction, error) {
	var tx Transaction
	var ratiosJson []byte
	err := r.DB.QueryRow(`SELECT id, block_id, payer, amount, 
       description, created_at, ratios FROM transactions WHERE id = $1`, id).
		Scan(&tx.ID, &tx.BlockID, &tx.Payer, &tx.Amount, &tx.Description, &tx.CreatedAt, &ratiosJson)
	if err != nil {
		return tx, err
	}
	_ = json.Unmarshal(ratiosJson, &tx.Ratios)
	return tx, nil
}

func (r *TransactionRepository) GetDetails(id string) (map[string]float64, error) {
	details := make(map[string]float64)
	rows, err := r.DB.Query(`SELECT member_id, amount FROM transaction_details WHERE transaction_id = $1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var mid string
		var amt float64
		if err := rows.Scan(&mid, &amt); err != nil {
			return nil, err
		}
		details[mid] = amt
	}
	return details, nil
}

func (r *TransactionRepository) Delete(id string) error {
	_, err := r.DB.Exec(`DELETE FROM transaction_details WHERE transaction_id = $1`, id)
	if err != nil {
		return err
	}
	_, err = r.DB.Exec(`DELETE FROM transactions WHERE id = $1`, id)
	return err
}

func (r *TransactionRepository) UpdateTransaction(payload UpdateTransactionPayload) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Lấy block_id của transaction hiện tại
	var blockID string
	err = tx.QueryRow(`SELECT block_id FROM transactions WHERE id = $1`, payload.ID).Scan(&blockID)
	if err != nil {
		return fmt.Errorf("failed to get block_id: %w", err)
	}

	var lock bool
	err = tx.QueryRow(`SELECT locked FROM blocks WHERE id = $1`, blockID).Scan(&lock)
	if err != nil || lock {
		return fmt.Errorf("failed to get block or block is locking: %w", err)
	}

	// Cập nhật transaction
	ratiosJSON, err := json.Marshal(payload.Ratios)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE transactions SET description=$1, amount=$2, payer=$3, ratios=$4 WHERE id=$5`,
		payload.Description, payload.Amount, payload.Payer, ratiosJSON, payload.ID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM transaction_details WHERE transaction_id=$1`, payload.ID)
	if err != nil {
		return err
	}

	var totalRatio float64
	for _, ratio := range payload.Ratios {
		totalRatio += ratio
	}
	if totalRatio == 0 {
		return fmt.Errorf("total ratio cannot be zero")
	}

	for memberID, ratio := range payload.Ratios {
		amount := payload.Amount * (ratio / totalRatio)
		_, err = tx.Exec(
			`INSERT INTO transaction_details (transaction_id, member_id, amount) VALUES ($1, $2, $3)`,
			payload.ID, memberID, amount,
		)
		if err != nil {
			return err
		}
	}

	if err := r.UpdateMembersDebtTx(tx, blockID); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *TransactionRepository) UpdateMembersDebtTx(tx *sql.Tx, blockID string) error {
	_, err := tx.Exec(`
		UPDATE members m
		SET debt = COALESCE((
			SELECT SUM(
				CASE 
					WHEN t.payer = m.id THEN t.amount
					ELSE 0
				END
			) - COALESCE((
				SELECT SUM(td.amount)
				FROM transaction_details td
				JOIN transactions t2 ON td.transaction_id = t2.id
				WHERE t2.block_id = m.block_id AND td.member_id = m.id
			), 0)
			FROM transactions t
			WHERE t.block_id = m.block_id
		), 0)
		WHERE m.block_id = $1`, blockID)
	return err
}
