package repository

import "database/sql"

type MemberRepository struct {
	DB *sql.DB
}

func NewMemberRepository(db *sql.DB) *MemberRepository {
	return &MemberRepository{DB: db}
}

func (r *MemberRepository) GetAll() ([]Member, error) {
	rows, err := r.DB.Query(`SELECT id, block_id, name, ratio, debt FROM members`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []Member
	for rows.Next() {
		var m Member
		if err := rows.Scan(&m.ID, &m.BlockID, &m.Name, &m.Ratio, &m.Debt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func (r *MemberRepository) GetByBlockID(blockID string) ([]Member, error) {
	rows, err := r.DB.Query(`SELECT id, block_id, name, ratio, debt FROM members WHERE block_id = $1`, blockID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []Member
	for rows.Next() {
		var m Member
		if err := rows.Scan(&m.ID, &m.BlockID, &m.Name, &m.Ratio, &m.Debt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, nil
}

func (r *MemberRepository) Create(members []Member) error {
	stmt, err := r.DB.Prepare(`INSERT INTO members (id, block_id, name, ratio, debt) VALUES ($1, $2, $3, $4, $5)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, m := range members {
		if _, err := stmt.Exec(m.ID, m.BlockID, m.Name, m.Ratio, m.Debt); err != nil {
			return err
		}
	}
	return nil
}

func (r *MemberRepository) UpdateDebt(id string, delta int) error {
	_, err := r.DB.Exec(`UPDATE members SET debt = debt + $1 WHERE id = $2`, delta, id)
	return err
}

func (r *MemberRepository) GetDebtsByBlockID(blockID string) (map[string]int, error) {
	rows, err := r.DB.Query(`SELECT name, debt FROM members WHERE block_id = $1`, blockID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[string]int{}
	for rows.Next() {
		var name string
		var debt int
		if err := rows.Scan(&name, &debt); err != nil {
			return nil, err
		}
		result[name] = debt
	}
	return result, nil
}
