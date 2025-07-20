// Filename: main.go
package main

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	_ "github.com/lib/pq"
	"log"
	authenhandler "my-source/sheet-payment/be/biz/auth"
	"my-source/sheet-payment/be/factory"

	_ "my-source/sheet-payment/be/docs"

	fiber "github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	_ "github.com/swaggo/files"
)

// @title Expense Tracker API
// @version 1.0
// @description API for managing shared expenses by month
// @host localhost:3000
// @BasePath /

// @contact.name Bui Phung Huu Duc
// @contact.email ducbph.x@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @Summary Get all transactions for a block
// @Tags transactions
// @Security BearerAuth
// @Produce json
// @Param month path string true "Month"
// @Success 200 {array} repository.Transaction
// @Router /blocks/{month}/transactions [get]
func getTransactionsByBlock(c *fiber.Ctx) error {
	return factory.GetBiz().GetTransactionsByBlock(c)
}

// @Summary Get all members
// @Tags members
// @Security BearerAuth
// @Produce json
// @Success 200 {array} repository.Member
// @Router /members [get]
func getAllMembers(c *fiber.Ctx) error {
	return factory.GetBiz().GetAllMembers(c)
}

// @Summary Lock a block
// @Tags blocks
// @Security BearerAuth
// @Param month path string true "Month"
// @Success 200 {string} string "locked"
// @Router /blocks/{month}/lock [post]
func lockBlock(c *fiber.Ctx) error {
	return factory.GetBiz().LockBlock(c)
}

// @Summary Unlock a block
// @Tags blocks
// @Security BearerAuth
// @Param month path string true "Month"
// @Success 200 {string} string "unlocked"
// @Router /blocks/{month}/unlock [post]
func unlockBlock(c *fiber.Ctx) error {
	return factory.GetBiz().UnlockBlock(c)
}

// @Summary Get summary of member debts in a block
// @Tags blocks
// @Security BearerAuth
// @Produce json
// @Param month path string true "Month"
// @Success 200 {object} map[string]int
// @Router /blocks/{month}/summary [get]
func getSummary(c *fiber.Ctx) error {
	return factory.GetBiz().GetSummary(c)
}

// @Summary Add a transaction to a block
// @Tags transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param month path string true "Month"
// @Param body body repository.Transaction true "Transaction info"
// @Success 200 {object} map[string]interface{}
// @Router /blocks/{month}/transactions [post]
func addTransaction(c *fiber.Ctx) error {
	return factory.GetBiz().AddTransaction(c)
}

// @Summary Create a new block
// @Tags blocks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body repository.CreateBlock true "Month and members"
// @Success 200 {object} repository.Block
// @Router /blocks [post]
func createBlock(c *fiber.Ctx) error {
	return factory.GetBiz().CreateBlock(c)
}

// @Summary Get members of a specific block
// @Tags members
// @Security BearerAuth
// @Produce json
// @Param month path string true "Month"
// @Success 200 {array} repository.Member
// @Router /blocks/{month}/members [get]
func getMembersByBlock(c *fiber.Ctx) error {
	return factory.GetBiz().GetMembersByLockId(c)
}

// @Summary Delete a transaction by ID
// @Description Removes a transaction and updates member debts accordingly
// @Tags transactions
// @Security BearerAuth
// @Param id path string true "Transaction ID"
// @Success 204 {string} string "No Content"
// @Router /transactions/{id} [delete]
func deleteTransaction(c *fiber.Ctx) error {
	return factory.GetBiz().DeleteTransaction(c)
}

// @Summary Login and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body authenhandler.LoginRequest true "Credentials"
// @Success 200 {object} map[string]string
// @Router /login [post]
func login(c *fiber.Ctx) error {
	return factory.GetAuth().Login(c)
}

// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param body body authenhandler.LoginRequest true "New user"
// @Success 200 {object} map[string]string
// @Router /register [post]
func register(c *fiber.Ctx) error {
	return factory.GetAuth().Register(c)
}

// @Summary Get all user logs
// @Description Retrieve all user logs
// @Tags logs
// @Security BearerAuth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} repository.UserLog
// @Router /logs [get]
func getLogs(c *fiber.Ctx) error {
	return factory.GetLogging().GetLogs(c)
}

// GetAllBlocks godoc
// @Summary Get all blocks
// @Description Get list of all blocks
// @Tags blocks
// @Security BearerAuth
// @Produce json
// @Success 200 {array} repository.Block
// @Failure 500 {object} object
// @Router /blocks [get]
func getBlocks(c *fiber.Ctx) error {
	return factory.GetBiz().GetAllBlocks(c)
}

// DeleteBlock godoc
// @Summary      Xóa block
// @Description  Xóa block theo ID, đồng thời xóa toàn bộ members và transactions liên quan
// @Tags         blocks
// @Security BearerAuth
// @Param        blockID   path      string  true  "ID của block"
// @Success      204       "Xóa thành công"
// @Failure      400       {object}  map[string]string  "Invalid ID"
// @Failure      500       {object}  map[string]string  "Internal server error"
// @Router       /blocks/{blockID} [delete]
// @Security     ApiKeyAuth
func deleteBlock(c *fiber.Ctx) error {
	return factory.GetBiz().DeleteBlock(c)
}

// UpdateTransaction godoc
// @Summary      Cập nhật giao dịch
// @Description  Cập nhật mô tả, số tiền, người trả và tỉ lệ chia của một giao dịch
// @Tags         transactions
// @Security     BearerAuth
// @Param        id         path      string  true  "ID của giao dịch"
// @Accept       json
// @Produce      json
//
// @Param body body repository.UpdateTransactionPayload true "Update transaction payload"
// @Success      200       {object}  map[string]string  "Transaction updated"
// @Failure      400       {object}  map[string]string  "Invalid input"
// @Failure      404       {object}  map[string]string  "Transaction not found"
// @Failure      500       {object}  map[string]string  "Internal server error"
// @Router       /transactions/{id} [put]
func updateTransaction(c *fiber.Ctx) error {
	return factory.GetBiz().UpdateTransaction(c)
}

func main() {
	factory.Factory()
	app := factory.GetApp()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Post("/login", login)
	app.Post("/register", register)

	protected := app.Group("/", jwtware.New(jwtware.Config{
		SigningKey: authenhandler.JwtSecret,
	}))

	protected.Use(factory.GetLogging().LogUserActivity())
	protected.Get("/blocks", getBlocks)
	protected.Post("/blocks", createBlock)
	protected.Delete("/blocks/:blockId/", deleteBlock)
	protected.Post("/blocks/:month/transactions", addTransaction)
	protected.Get("/blocks/:month/transactions", getTransactionsByBlock)
	protected.Get("/blocks/:month/summary", getSummary)
	protected.Get("/members", getAllMembers)
	protected.Post("/blocks/:month/lock", lockBlock)
	protected.Post("/blocks/:month/unlock", unlockBlock)
	protected.Get("/blocks/:month/members", getMembersByBlock)
	protected.Delete("/transactions/:id", deleteTransaction)
	protected.Get("/logs", getLogs)
	protected.Put("/transactions/:id", updateTransaction)

	log.Fatal(app.Listen(":3000"))
}
