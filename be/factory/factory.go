package factory

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"log"
	mainbiz "my-source/sheet-payment/be/biz"
	authenhandler "my-source/sheet-payment/be/biz/auth"
	"my-source/sheet-payment/be/repository"
	"os"
)

var (
	db       *sql.DB
	app      *fiber.App
	bizInst  *mainbiz.MainBusiness
	authInst *authenhandler.AuthHandler
)

func initDB() {
	var err error
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=pgdb port=5432 user=postgres password=yourpassword dbname=expenses sslmode=disable"
	}
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	queries := []string{
		`CREATE TABLE IF NOT EXISTS blocks (
			id TEXT PRIMARY KEY,
			month TEXT UNIQUE,
			locked BOOLEAN
		)`,
		`CREATE TABLE IF NOT EXISTS members (
			id TEXT PRIMARY KEY,
			block_id TEXT,
			name TEXT,
			ratio FLOAT,
			debt INT,
			FOREIGN KEY (block_id) REFERENCES blocks(id)
		)`,
		`CREATE TABLE IF NOT EXISTS transactions (
			id TEXT PRIMARY KEY,
			block_id TEXT,
			payer TEXT,
			amount INT,
			description TEXT,
			created_at TIMESTAMP,
			ratios JSONB,
			FOREIGN KEY (block_id) REFERENCES blocks(id)
		)`,
		`CREATE TABLE IF NOT EXISTS transaction_details (
			transaction_id TEXT,
			member_id TEXT,
			amount INT,
			PRIMARY KEY (transaction_id, member_id)
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
		)`,
	}
	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			log.Fatal(err)
		}
	}
}

func Factory() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it")
	}

	initDB()

	memberRepo := repository.NewMemberRepository(db)
	blockRepo := repository.NewBlockRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	userRepo := repository.NewUserRepository(db)
	authInst = authenhandler.NewAuthHandler(userRepo)
	bizInst = mainbiz.NewMainBusiness(memberRepo, blockRepo, transactionRepo)
	app = fiber.New()
}

func GetApp() *fiber.App {
	return app
}

func GetBiz() *mainbiz.MainBusiness {
	return bizInst
}

func GetAuth() *authenhandler.AuthHandler {
	return authInst
}
