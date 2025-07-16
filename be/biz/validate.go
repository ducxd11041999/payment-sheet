package mainbiz

import (
	"github.com/gofiber/fiber/v2"
)

func (mb *MainBusiness) ValidateMemberInMonth(month string, member map[string]float64, payerId string) error {
	blockId, locked, err := mb.blockRepo.GetIDByMonth(month)
	if err != nil {
		return err
	}

	if locked {
		return fiber.ErrForbidden
	}

	memberInBlock, err := mb.memberRepo.GetByBlockID(blockId)
	mp := map[string]bool{}
	for _, m := range memberInBlock {
		mp[m.ID] = true
	}

	for k, _ := range member {
		if !mp[k] {
			return fiber.ErrNotFound
		}
	}

	if _, ok := member[payerId]; !ok {
		return fiber.ErrForbidden
	}

	return nil
}
