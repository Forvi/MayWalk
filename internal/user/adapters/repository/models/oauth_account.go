package models

import "time"

const YandexProvider string = "yandex"

type OAuthAccount struct {
	ID             string    `db:"id"`
	UserID         string    `db:"user_id"`
	Provider       string    `db:"provider"`
	ProviderUserID string    `db:"provider_user_id"`
	CreatedAt      time.Time `db:"created_at"`
}

type OAuthAccountCreate struct {
	UserID         string `json:"user_id"`
	Provider       string `json:"provider"`
	ProviderUserID string `json:"provider_user_id"`
}
