package api

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"net/url"
	"strconv"
	"time"

	"crdx.org/dex/cmd/dexd/env"
	"crdx.org/dex/db"
	"crdx.org/dex/pkg/types"
	"github.com/gofiber/fiber/v3"
	"github.com/samber/lo"
)

func Expose(c fiber.Ctx) error {
	var req types.ExposeRequest
	lo.Must0(c.Bind().Body(&req))

	item, found := db.FindItemByRef(req.Ref, req.Ref)
	if !found {
		return errRefNotFound(c, req.Ref)
	}

	var expiresAt sql.Null[time.Time]
	if req.ExpirySeconds > 0 {
		expiresAt = sql.Null[time.Time]{
			V:     time.Now().Add(time.Duration(req.ExpirySeconds) * time.Second),
			Valid: true,
		}
	} else {
		existingTokens := db.FindTokensByItemID(item.ID)
		for _, t := range existingTokens {
			if !t.ExpiresAt.Valid {
				deployURL, _ := url.JoinPath(env.BaseURL(), "deploy", t.Token)
				return c.JSON(types.ExposeResponse{
					PublicURL: item.URL(),
					DeployURL: deployURL,
					Warning:   "item already has a non-expiring token",
				})
			}
		}
	}

	token := generateToken()

	db.CreateToken(&db.Token{
		ItemID:    item.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	})

	deployURL, _ := url.JoinPath(env.BaseURL(), "deploy", token)

	res := types.ExposeResponse{
		PublicURL: item.URL(),
		DeployURL: deployURL,
	}

	if expiresAt.Valid {
		res.ExpiresAt = expiresAt.V.Format(time.RFC3339)
	}

	return c.JSON(res)
}

func generateToken() string {
	bytes := make([]byte, 32)
	lo.Must1(rand.Read(bytes))
	return base64.RawURLEncoding.EncodeToString(bytes)
}

func DeployURLs(c fiber.Ctx) error {
	tokens := db.FindTokens()

	items := make([]types.DeployURLsResponseItem, 0, len(tokens))
	for _, token := range tokens {
		item, found := db.FindItem(token.ItemID)
		if !found {
			continue
		}

		deployURL, _ := url.JoinPath(env.BaseURL(), "deploy", token.Token)

		responseItem := types.DeployURLsResponseItem{
			ID:        token.ID,
			Ref:       item.Ref(),
			PublicURL: item.URL(),
			DeployURL: deployURL,
			CreatedAt: token.CreatedAt.Format(time.RFC3339),
		}

		if token.ExpiresAt.Valid {
			responseItem.ExpiresAt = token.ExpiresAt.V.Format(time.RFC3339)
		}

		items = append(items, responseItem)
	}

	return c.JSON(types.DeployURLsResponse{Tokens: items})
}

func DeployURLsDelete(c fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return FailureResponse(c, 400, "invalid id")
	}

	token, found := db.FindToken(id)
	if !found {
		return FailureResponse(c, 404, "deploy url not found")
	}

	token.Delete()

	return c.JSON(types.DeployURLsDeleteResponse{Message: "deleted"})
}

func Deployments(c fiber.Ctx) error {
	deployments := db.FindDeployments()

	items := make([]types.DeploymentsResponseItem, 0, len(deployments))
	for _, d := range deployments {
		token, found := db.FindTokenUnscoped(d.TokenID)
		if !found {
			continue
		}

		item, found := db.FindItemUnscoped(token.ItemID)
		if !found {
			continue
		}

		items = append(items, types.DeploymentsResponseItem{
			ID:        d.ID,
			Ref:       item.Ref(),
			Change:    d.Note,
			Deployer:  d.Deployer,
			IPAddress: d.IPAddress,
			UserAgent: d.UserAgent,
			CreatedAt: d.CreatedAt.Format(time.RFC3339),
			Deleted:   item.DeletedAt.Valid,
		})
	}

	if c.Query("all") != "1" && len(items) > 10 {
		items = items[len(items)-10:]
	}

	return c.JSON(types.DeploymentsResponse{Deployments: items})
}
