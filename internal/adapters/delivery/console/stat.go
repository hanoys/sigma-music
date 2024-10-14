package console

import (
	"context"
	"fmt"
	"github.com/hanoys/sigma-music/internal/adapters/delivery/console/dto"
)

func (h *Handler) Listen(c *Console) {
	err := h.verifyUserAuth(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	id, err := readID()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = h.statService.Add(context.Background(), c.UserID, id)
	if err != nil {
		fmt.Println("listen error")
	}
}

func (h *Handler) GetStat(c *Console) {
	err := h.verifyUserAuth(c)
	if err != nil {
		fmt.Println(err)
		return
	}

	rep, err := h.statService.FormReport(context.Background(), c.UserID)
	if err != nil {
		fmt.Println(err)
	}

	dto.NewListenReportDTO(rep).Print()
}
