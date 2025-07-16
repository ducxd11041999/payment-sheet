package repository

import "time"

type Member struct {
	ID      string  `json:"id"`
	BlockID string  `json:"block_id"`
	Name    string  `json:"name"`
	Ratio   float64 `json:"ratio"`
	Debt    int     `json:"debt"`
}

type Transaction struct {
	ID          string             `json:"id"`
	BlockID     string             `json:"block_id"`
	Description string             `json:"description"`
	Amount      int                `json:"amount"`
	Payer       string             `json:"payer"`
	Details     map[string]int     `json:"details"`
	Ratios      map[string]float64 `json:"ratios"`
	CreatedAt   time.Time          `json:"created_at"`
}

type Block struct {
	ID           string         `json:"id"`
	Month        string         `json:"month"`
	Locked       bool           `json:"locked"`
	Members      []*Member      `json:"members"`
	Transactions []*Transaction `json:"transactions"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"` // Hashed password
}

type CreateBlock struct {
	Month   string    `json:"month"`
	Members []*Member `json:"members"`
}
