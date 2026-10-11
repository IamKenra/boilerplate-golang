package message

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestReturnBadRequest_RedactsDeveloperMessageInRelease(t *testing.T) {
	t.Setenv("APP_ENV", "release")

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		ReturnBadRequest(c, "debug detail")
		return nil
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body["message for developer"] != nil {
		t.Fatalf("expected developer message to be nil, got: %v", body["message for developer"])
	}
}

func TestReturnBadRequest_ExposesDeveloperMessageInDevelopment(t *testing.T) {
	t.Setenv("APP_ENV", "development")

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		ReturnBadRequest(c, "debug detail")
		return nil
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test failed: %v", err)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	got, ok := body["message for developer"].(string)
	if !ok {
		t.Fatalf("expected developer message as string, got: %T", body["message for developer"])
	}
	if got != "debug detail" {
		t.Fatalf("unexpected developer message: %s", got)
	}
}
