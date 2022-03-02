package routes

import (
	"time"

	"github.com/ArsalanKm/url-shortner/helpers"
	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL            string `json:"url"`
	CustomShort    string `json:"custom_short"`
	Expiry         string `json:"expiry"`
	XrateRemaining string `json:"rate_remaining"`
	XrateLimitRest string `json:"rate_limit_rest"`
}

func ShortenURL(ctx *fiber.Ctx) error {
	body := new(request)
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Can not parse json"})
	}
	// implement rate limiting

	// check if the input url is actual url
	if !govalidator.IsURL(body.URL) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad url"})
	}

	// check for domain error
	if !helpers.RemoveDomainError(body.URL) {
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "service unavailable"})
	}

	// enforce https,SSl

	body.URL = helpers.EnforceHttp(body.URL)
	return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "service unavailable"})
}
