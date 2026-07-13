package api

import (
	"fmt"
	"net/http"
	"strings"

	"crdx.org/dex/pkg/types"
	"github.com/gofiber/fiber/v3"
)

const reservedLabelPrefix = "-/"

func isReservedLabel(label string) bool {
	return strings.HasPrefix(label, reservedLabelPrefix)
}

func errRefNotFound(c fiber.Ctx, ref string) error {
	return FailureResponse(c, http.StatusNotFound, "ref %s not found", ref)
}

func errDuplicateRef(c fiber.Ctx, ref string) error {
	return FailureResponse(c, http.StatusConflict, "duplicate ref %s", ref)
}

func errReservedRef(c fiber.Ctx, ref string) error {
	return FailureResponse(c, http.StatusBadRequest, "reserved ref %s", ref)
}

func errInvalidKind(c fiber.Ctx, kind string) error {
	return FailureResponse(c, http.StatusBadRequest, "invalid kind %s", kind)
}

func errInvalidContentHash(c fiber.Ctx) error {
	return FailureResponse(c, http.StatusBadRequest, "invalid content hash")
}

func errUnauthorized(c fiber.Ctx, msg string) error {
	return FailureResponse(c, http.StatusUnauthorized, "%s", msg)
}

func FailureResponse(c fiber.Ctx, status int, format string, a ...any) error {
	return c.
		Status(status).
		JSON(types.FailureResponse{
			Message: fmt.Sprintf(format, a...),
		})
}
