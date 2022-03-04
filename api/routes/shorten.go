package routes

import (
	"os"
	"strconv"
	"time"

	"github.com/ArsalanKm/url-shortner/database"
	"github.com/ArsalanKm/url-shortner/helpers"
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	r2 := database.CreateClient(1)
	defer r2.Close()
	val, err := r2.Get(database.Ctx, ctx.IP()).Result()
	if err == redis.Nil {
		_ = r2.Set(database.Ctx, ctx.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		val, _ = r2.Get(database.Ctx, ctx.IP()).Result()
		valint, _ := strconv.Atoi(val)
		if valint <= 0 {
			limit, _ := r2.TTL(database.Ctx, ctx.IP()).Result()
			return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "rate limit exceeded", "rate_limit_reset": limit / time.Nanosecond / time.Minute})
		}
	}

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
	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}
	r := database.CreateClient(0)
	defer r.Close()

	val, _ = r.Get(database.Ctx, id).Result()

	if val != "" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Url custom short is already been used"})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}
	err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Unable to connect to server"})
	}
	r2.Decr(database.Ctx, ctx.IP())
	return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "service unavailable"})
}
