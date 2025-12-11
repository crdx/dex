package types

import (
	"database/sql"
	"time"
)

type SetRequest struct {
	Ref   string `json:"ref"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SetResponse struct {
	Modified bool `json:"modified"`
}

type CopyRequest struct {
	From string `json:"from_ref"`
	To   string `json:"to_label"`
}

type CopyResponse struct {
	URL string `json:"url"`
}

type MoveRequest struct {
	FromRef string `json:"from_ref"`
	ToLabel string `json:"to_label"`
}

type MoveResponse struct {
	URL string `json:"url"`
}

type UploadRequest struct {
	Label       string `json:"label"`
	Kind        string `json:"kind"`
	Content     string `json:"content"`
	ContentHash string `json:"content_hash"`
	ContentType string `json:"content_type"`
	Force       bool   `json:"force"`
}

type UploadResponse struct {
	URL string `json:"url"`
}

type DeleteResponse struct {
	Message string `json:"message"`
}

type GCResponse struct {
	DeletedHashes []string `json:"deleted"`
}

type CatResponse struct {
	Kind        string `json:"kind"`
	Content     []byte `json:"content"`
	ContentType string `json:"content_type"`
	IsText      bool   `json:"is_text"`
}

type TailResponse struct {
	Logs []TailResponseItem `json:"logs"`
}

type TailResponseItem struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Method    string    `json:"method"`
	Request   string    `json:"request"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	BaseURL   string    `json:"base_url"`
	Renamed   bool      `json:"renamed"`
	Exists    bool      `json:"exists"`
	Ref       string    `json:"ref"`
}

type ListResponse struct {
	Items []ListResponseItem `json:"items"`
}

type ListResponseItem struct {
	Ref         string              `json:"ref"`
	URL         string              `json:"url"`
	Kind        string              `json:"kind"`
	ContentType string              `json:"mime"`
	IsText      bool                `json:"is_text"`
	Hits        int64               `json:"hits"`
	Hash        string              `json:"hash"`
	Meta        string              `json:"meta"`
	LastHitAt   sql.Null[time.Time] `json:"last_hit_at"`
}

type FailureResponse struct {
	Message string `json:"message"`
}

type ExposeRequest struct {
	Ref           string `json:"ref"`
	ExpirySeconds int64  `json:"expiry_seconds,omitempty"`
}

type ExposeResponse struct {
	PublicURL string `json:"public_url"`
	DeployURL string `json:"deploy_url"`
	ExpiresAt string `json:"expires_at,omitempty"`
	Warning   string `json:"warning,omitempty"`
}

type DeployURLsResponse struct {
	Tokens []DeployURLsResponseItem `json:"tokens"`
}

type DeployURLsResponseItem struct {
	ID        int64  `json:"id"`
	Ref       string `json:"ref"`
	PublicURL string `json:"public_url"`
	DeployURL string `json:"deploy_url"`
	ExpiresAt string `json:"expires_at,omitempty"`
	CreatedAt string `json:"created_at"`
}

type DeployURLsDeleteResponse struct {
	Message string `json:"message"`
}

type DeploymentsResponse struct {
	Deployments []DeploymentsResponseItem `json:"deployments"`
}

type DeploymentsResponseItem struct {
	ID        int64  `json:"id"`
	Ref       string `json:"ref"`
	Note      string `json:"note"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	CreatedAt string `json:"created_at"`
	Deleted   bool   `json:"deleted"`
}
