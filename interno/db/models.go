// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package db

import (
	"time"

	"github.com/google/uuid"
)

type Usuario struct {
	ID        uuid.UUID
	Email     string
	Nome      string
	Telefone  string
	Senha     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
