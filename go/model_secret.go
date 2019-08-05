/*
 * Secret Server
 *
 * This is an API of a secret service. You can save your secret by using the API. You can restrict the access of a secret after the certen number of views or after a certen period of time.
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

import (
	"time"
)

type Secret struct {

	// Unique hash to identify the secrets
	Hash string `json:"hash,omitempty"`

	// The secret itself
	SecretText string `json:"secretText,omitempty"`

	// The date and time of the creation
	CreatedAt time.Time `json:"createdAt,omitempty"`

	// The secret cannot be reached after this time
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`

	// How many times the secret can be viewed
	RemainingViews int32 `json:"remainingViews,omitempty"`
}
