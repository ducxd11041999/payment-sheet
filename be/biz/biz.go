package mainbiz

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"my-source/sheet-payment/be/repository"
	"time"
)

type MainBusiness struct {
	memberRepo      repository.IMemberRepository
	blockRepo       repository.IBlockRepository
	transactionRepo repository.ITransactionRepository
}

func NewMainBusiness(mrb repository.IMemberRepository, brp repository.IBlockRepository,
	trp repository.ITransactionRepository) *MainBusiness {
	return &MainBusiness{
		memberRepo:      mrb,
		blockRepo:       brp,
		transactionRepo: trp,
	}
}

func (mb *MainBusiness) GetAllMembers(c *fiber.Ctx) error {
	members, err := mb.memberRepo.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}

	return c.JSON(members)
}

func (mb *MainBusiness) GetMembersByLockId(c *fiber.Ctx) error {
	month := c.Params("month")
	blockID, _, err := mb.blockRepo.GetIDByMonth(month)
	if err != nil {
		return err
	}

	members, err := mb.memberRepo.GetByBlockID(blockID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}

	return c.JSON(members)
}

func (mb *MainBusiness) CreateBlock(c *fiber.Ctx) error {
	type Req struct {
		Month   string               `json:"month"`
		Members []*repository.Member `json:"members"`
	}

	var req Req
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	id := uuid.New().String()
	block := repository.Block{
		ID:      id,
		Month:   req.Month,
		Locked:  false,
		Members: req.Members,
	}

	err := mb.blockRepo.Create(block)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}

	return c.JSON(block)
}

func (mb *MainBusiness) LockBlock(c *fiber.Ctx) error {
	month := c.Params("month")
	err := mb.blockRepo.Lock(month)
	if err != nil {
		return err
	}
	return c.SendString("locked")
}

func (mb *MainBusiness) UnlockBlock(c *fiber.Ctx) error {
	month := c.Params("month")
	err := mb.blockRepo.Unlock(month)
	if err != nil {
		return err
	}
	return c.SendString("unlocked")
}

func (mb *MainBusiness) AddTransaction(c *fiber.Ctx) error {
	month := c.Params("month")
	blockId, isLock, err := mb.blockRepo.GetIDByMonth(month)
	if err != nil {
		return err
	}

	if isLock {
		return fiber.NewError(fiber.StatusForbidden, "Page is block")
	}

	type Req struct {
		Amount      float64            `json:"amount"`
		Description string             `json:"description"`
		Payer       string             `json:"payer"`
		Ratios      map[string]float64 `json:"ratios"`
	}

	var req Req
	if err := c.BodyParser(&req); err != nil {
		return err
	}

	err = mb.ValidateMemberInMonth(month, req.Ratios, req.Payer)
	if err != nil {
		return err
	}

	totalWeight := 0.0
	for _, w := range req.Ratios {
		totalWeight += w
	}
	if totalWeight == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "participants weight must be > 0")
	}

	txID := uuid.New().String()
	created := time.Now()
	tx := repository.Transaction{
		ID:          txID,
		BlockID:     blockId,
		Description: req.Description,
		Amount:      req.Amount,
		Payer:       req.Payer,
		Ratios:      req.Ratios,
		CreatedAt:   created,
	}

	if err := mb.transactionRepo.Add(tx); err != nil {
		return err
	}

	// Prepare details
	details := make(map[string]float64)
	for memberID, weight := range req.Ratios {
		share := req.Amount * (weight / totalWeight)
		details[memberID] = share
	}

	if err := mb.transactionRepo.AddDetails(txID, details); err != nil {
		return err
	}

	// Update debts
	for memberID, share := range details {
		if memberID == req.Payer {
			_ = mb.memberRepo.UpdateDebt(memberID, req.Amount-share)
		} else {
			_ = mb.memberRepo.UpdateDebt(memberID, -share)
		}
	}

	return c.JSON(fiber.Map{"id": txID, "block_id": blockId, "created_at": created})
}

func (mb *MainBusiness) GetSummary(c *fiber.Ctx) error {
	month := c.Params("month")
	blockID, _, err := mb.blockRepo.GetIDByMonth(month)
	if err != nil {
		return err
	}
	summary, err := mb.memberRepo.GetDebtsByBlockID(blockID)
	if err != nil {
		return err
	}

	return c.JSON(summary)
}

func (mb *MainBusiness) GetTransactionsByBlock(c *fiber.Ctx) error {
	month := c.Params("month")
	blockID, _, err := mb.blockRepo.GetIDByMonth(month)
	if err != nil {
		return nil
	}
	txs, err := mb.transactionRepo.GetByBlockID(blockID)
	return c.JSON(txs)
}

func (mb *MainBusiness) DeleteTransaction(c *fiber.Ctx) error {
	id := c.Params("id")
	tx, err := mb.transactionRepo.GetByID(id)
	if err != nil {
		return fiber.NewError(fiber.StatusForbidden, "not found tx")
	}

	_, err = mb.transactionRepo.GetDetails(id)
	if err != nil {
		return fiber.NewError(fiber.StatusForbidden, "not found tx details")
	}

	_, lock, er := mb.blockRepo.Get(tx.BlockID)
	if er != nil {
		return fiber.NewError(fiber.StatusForbidden, "not found block")
	}

	if lock {
		return fiber.NewError(fiber.StatusForbidden, "locked by this block")
	}

	// Reverse debts
	totalWeight := 0.0
	for _, w := range tx.Ratios {
		totalWeight += w
	}

	for memberID, weight := range tx.Ratios {
		share := tx.Amount * (weight / totalWeight)
		if memberID == tx.Payer {
			// Payer previously gained (credit), now subtract it back
			err = mb.memberRepo.UpdateDebt(memberID, -(tx.Amount - share))
		} else {
			// Participant paid less before, now cancel the minus
			err = mb.memberRepo.UpdateDebt(memberID, share)
		}
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	return mb.transactionRepo.Delete(id)
}

func (mb *MainBusiness) GetAllBlocks(c *fiber.Ctx) error {
	blocks, err := mb.blockRepo.GetAllBlocks()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get blocks",
		})
	}
	return c.JSON(blocks)
}

func (mb *MainBusiness) DeleteBlock(c *fiber.Ctx) error {
	blockID := c.Params("blockID")
	err := mb.blockRepo.DeleteBlock(blockID)
	if err != nil {
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	}

	return nil
}

func (mb *MainBusiness) UpdateTransaction(c *fiber.Ctx) error {
	id := c.Params("id")
	var body struct {
		Description string             `json:"description"`
		Amount      float64            `json:"amount"`
		Payer       string             `json:"payer"`
		Ratios      map[string]float64 `json:"ratios"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	payload := repository.UpdateTransactionPayload{
		ID:          id,
		Description: body.Description,
		Amount:      body.Amount,
		Payer:       body.Payer,
		Ratios:      body.Ratios,
	}

	if err := mb.transactionRepo.UpdateTransaction(payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Transaction updated successfully"})
}
