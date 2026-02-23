package http

import (
	"boilerplate-golang/utils/message"

	"github.com/gofiber/fiber/v2"
)

func (handler *HttpMasterdataHandler) Health(ctx *fiber.Ctx) error {
	status, err := handler.masterdataUsecase.Health()
	if err != nil {
		message.ReturnInternalServerError(ctx, err.Error())
		return nil
	}

	message.ReturnOkNoParam(ctx, status)
	return nil
}
