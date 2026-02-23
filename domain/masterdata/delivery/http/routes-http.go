package http

import (
	"github.com/gofiber/fiber/v2"

	"boilerplate-golang/domain/masterdata"
	"boilerplate-golang/utils/middleware"
)

type HttpMasterdataHandler struct {
	masterdataUsecase masterdata.Usecase
}

func NewHttpMasterdataHandler(masterdataUC masterdata.Usecase, httpRouter *fiber.App) {
	handler := &HttpMasterdataHandler{
		masterdataUsecase: masterdataUC,
	}

	//public route
	public := httpRouter.Group("/api/v1/masterdata")
	public.Use(middleware.PublicMiddleware)

	//health check route
	public.Get("/health", handler.Health)
}
