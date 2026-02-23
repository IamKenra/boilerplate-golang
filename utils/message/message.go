package message

import (
	"log"
	"net/http"
	"os"
	"reflect"
	"boilerplate-golang/utils/slices"

	"github.com/gofiber/fiber/v2"
)

const PUBLIC = "public"
const RELEASE = "release"
const ERROR_STATUS = "Layanan sedang mengalami gangguan"
const BAD_REQUEST_STATUS = "Input tidak valid"
const SUCCESS_UPDATE_STATUS = "Berhasil melakukan pembaharuan data"
const SUCCESS_INSERT_STATUS = "Berhasil melakukan penambahan data"

func measureLength(data interface{}) int {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Slice:
		return v.Len()
	case reflect.Struct:
		return 1
	}
	return 0
}

func castNilSliceToArray(data interface{}) interface{} {
	if reflect.ValueOf(data).Kind() == reflect.Slice {
		if reflect.ValueOf(data).Len() < 1 {
			return []interface{}{}
		}
	}
	return data
}

func ReturnOk(ctx *fiber.Ctx, data interface{}, param interface{}) {
	requestSource, ok := ctx.Locals("requestSource").(string)
	if !ok || requestSource == "" {
		ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Request source not found"})
		return
	}
	if requestSource == PUBLIC {
		slices.RemoveIDFields(&data)
	}
	if os.Getenv("GIN_MODE") == RELEASE {
		param = nil
	}
	ctx.Status(http.StatusOK).JSON(fiber.Map{
		"param": param,
		"count": measureLength(data),
		"data":  castNilSliceToArray(data),
	})
}

func ReturnOkNoParam(ctx *fiber.Ctx, data interface{}) {
	requestSource, ok := ctx.Locals("requestSource").(string)
	if !ok || requestSource == "" {
		ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Request source not found"})
		return
	}
	if requestSource == PUBLIC {
		slices.RemoveIDFields(&data)
	}
	ctx.Status(http.StatusOK).JSON(fiber.Map{
		"count": measureLength(data),
		"data":  castNilSliceToArray(data),
	})
}

func ReturnPagination(ctx *fiber.Ctx, data interface{}, param interface{}, count int) {
	requestSource, ok := ctx.Locals("requestSource").(string)
	if !ok || requestSource == "" {
		ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Request source not found"})
		return
	}
	if requestSource == PUBLIC {
		slices.RemoveIDFields(&data)
	}
	ctx.Status(http.StatusOK).JSON(fiber.Map{
		"param": param,
		"count": count,
		"data":  castNilSliceToArray(data),
	})
}

func ReturnBadRequest(ctx *fiber.Ctx, errorMessage interface{}) {
	rawBody := ctx.Body()
	if len(rawBody) == 0 {
		log.Print("error: ", ctx.OriginalURL(), " - ", errorMessage)
	} else {
		log.Print("error: ", ctx.OriginalURL(), " - ", errorMessage, " - ", string(rawBody))
	}
	if os.Getenv("GIN_MODE") == RELEASE {
		errorMessage = nil
	}
	ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
		"message":               BAD_REQUEST_STATUS,
		"message for developer": errorMessage,
	})
}

func ReturnInternalServerError(ctx *fiber.Ctx, errorMessage ...string) {
	rawBody := ctx.Body()
	if len(rawBody) == 0 {
		log.Print("error: ", ctx.OriginalURL(), " - ", errorMessage)
	} else {
		log.Print("error: ", ctx.OriginalURL(), " - ", errorMessage, " - ", string(rawBody))
	}
	if os.Getenv("GIN_MODE") == RELEASE {
		errorMessage = nil
	}
	ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
		"message":               ERROR_STATUS,
		"message for developer": errorMessage,
	})
}

func ReturnSuccessInsert(ctx *fiber.Ctx, data interface{}) {
	ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"result": true,
		"data":   data,
		"info":   SUCCESS_INSERT_STATUS,
	})
}

func ReturnSuccessDelete(ctx *fiber.Ctx) {
	ctx.SendStatus(fiber.StatusNoContent)
}

func ReturnSuccessUpdate(ctx *fiber.Ctx) {
	ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": SUCCESS_UPDATE_STATUS,
	})
}

func ReturnNotFound(ctx *fiber.Ctx, message string) {
	if os.Getenv("GIN_MODE") == RELEASE {
		message = ""
	}
	ctx.Status(http.StatusNotFound).JSON(fiber.Map{
		"message":               "Data tidak ditemukan.",
		"message for developer": message,
	})
}

func ReturnSuccessInsertSimple(ctx *fiber.Ctx) {
	ctx.Status(http.StatusCreated).JSON(fiber.Map{
		"result":  "success",
		"message": SUCCESS_INSERT_STATUS,
	})
}

func ReturnSuccessUpdateSimple(ctx *fiber.Ctx) {
	ctx.Status(http.StatusOK).JSON(fiber.Map{
		"result":  "success",
		"message": SUCCESS_UPDATE_STATUS,
	})
}
