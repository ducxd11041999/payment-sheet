package repository

type IBlockRepository interface {
	GetAllBlocks() ([]Block, error)
	Get(id string) (string, bool, error)
	GetIDByMonth(month string) (string, bool, error)
	Lock(month string) error
	Unlock(month string) error
	Create(block Block) error
	DeleteBlock(blockID string) error
}

type IMemberRepository interface {
	GetAll() ([]Member, error)
	GetByBlockID(blockID string) ([]Member, error)
	Create(members []Member) error
	UpdateDebt(id string, delta float64) error
	GetDebtsByBlockID(blockID string) (map[string]int, error)
}

type ITransactionRepository interface {
	GetByID(id string) (Transaction, error)
	GetDetails(id string) (map[string]float64, error)
	GetByBlockID(blockID string) ([]Transaction, error)
	Add(tx Transaction) error
	AddDetails(txID string, details map[string]float64) error
	Delete(id string) error
	UpdateTransaction(payload UpdateTransactionPayload) error
}

type IUserRepository interface {
	GetByUsername(username string) (*User, error)
	Create(user *User) error
}

type ILogging interface {
	Write(logEntry UserLog) error
	GetAllLogs() ([]UserLog, error)
}
