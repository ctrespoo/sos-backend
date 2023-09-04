// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: query.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const criarCategoria = `-- name: CriarCategoria :one
INSERT INTO "categorias" ("nome", "imagem")
VALUES ($1, $2)
RETURNING "id"
`

type CriarCategoriaParams struct {
	Nome   string
	Imagem string
}

func (q *Queries) CriarCategoria(ctx context.Context, arg CriarCategoriaParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, criarCategoria, arg.Nome, arg.Imagem)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const criarProduto = `-- name: CriarProduto :one
INSERT INTO "produtos" (
        "nome",
        "descricao",
        "preco",
        "unidade_medida",
        "quantidade_pacote",
        "peso",
        "ativo",
        "ordem",
        "imagem"
    )
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING "id"
`

type CriarProdutoParams struct {
	Nome             string
	Descricao        string
	Preco            pgtype.Numeric
	UnidadeMedida    TipoUnidadeMedida
	QuantidadePacote int32
	Peso             pgtype.Numeric
	Ativo            bool
	Ordem            int32
	Imagem           string
}

func (q *Queries) CriarProduto(ctx context.Context, arg CriarProdutoParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, criarProduto,
		arg.Nome,
		arg.Descricao,
		arg.Preco,
		arg.UnidadeMedida,
		arg.QuantidadePacote,
		arg.Peso,
		arg.Ativo,
		arg.Ordem,
		arg.Imagem,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const criarUsuario = `-- name: CriarUsuario :one
INSERT INTO usuarios (
        nome,
        email,
        telefone,
        ativo
    )
VALUES (
        $1,
        $2,
        $3,
        $4
    )
RETURNING id, email, nome, telefone, role, ativo, created_at, updated_at
`

type CriarUsuarioParams struct {
	Nome     string
	Email    string
	Telefone string
	Ativo    bool
}

func (q *Queries) CriarUsuario(ctx context.Context, arg CriarUsuarioParams) (Usuario, error) {
	row := q.db.QueryRow(ctx, criarUsuario,
		arg.Nome,
		arg.Email,
		arg.Telefone,
		arg.Ativo,
	)
	var i Usuario
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Nome,
		&i.Telefone,
		&i.Role,
		&i.Ativo,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const pegarTodasCategorias = `-- name: PegarTodasCategorias :many
SELECT id, nome, imagem, created_at, updated_at
FROM "categorias"
ORDER BY "updated_at" DESC
LIMIT $1 OFFSET $2
`

type PegarTodasCategoriasParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) PegarTodasCategorias(ctx context.Context, arg PegarTodasCategoriasParams) ([]Categoria, error) {
	rows, err := q.db.Query(ctx, pegarTodasCategorias, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Categoria
	for rows.Next() {
		var i Categoria
		if err := rows.Scan(
			&i.ID,
			&i.Nome,
			&i.Imagem,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const pegarTodosProdutos = `-- name: PegarTodosProdutos :many
SELECT p.id AS produto_id,
    p.nome AS produto_nome,
    p.descricao AS produto_descricao,
    p.preco,
    p.unidade_medida,
    p.quantidade_pacote,
    p.peso,
    p.ativo AS produto_ativo,
    p.ordem,
    p.imagem AS produto_imagem,
    array_to_string(array_agg(c.nome), ',') AS categorias
FROM produtos p
    LEFT JOIN "_CategoriaToProduto" cp ON p.id = cp."B"
    LEFT JOIN categorias c ON cp."A" = c.id
GROUP BY p.id,
    p.nome,
    p.descricao,
    p.preco,
    p.unidade_medida,
    p.quantidade_pacote,
    p.peso,
    p.ativo,
    p.ordem,
    p.imagem
ORDER BY p.updated_at DESC
LIMIT $1 OFFSET $2
`

type PegarTodosProdutosParams struct {
	Limit  int32
	Offset int32
}

type PegarTodosProdutosRow struct {
	ProdutoID        uuid.UUID
	ProdutoNome      string
	ProdutoDescricao string
	Preco            pgtype.Numeric
	UnidadeMedida    TipoUnidadeMedida
	QuantidadePacote int32
	Peso             pgtype.Numeric
	ProdutoAtivo     bool
	Ordem            int32
	ProdutoImagem    string
	Categorias       string
}

func (q *Queries) PegarTodosProdutos(ctx context.Context, arg PegarTodosProdutosParams) ([]PegarTodosProdutosRow, error) {
	rows, err := q.db.Query(ctx, pegarTodosProdutos, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PegarTodosProdutosRow
	for rows.Next() {
		var i PegarTodosProdutosRow
		if err := rows.Scan(
			&i.ProdutoID,
			&i.ProdutoNome,
			&i.ProdutoDescricao,
			&i.Preco,
			&i.UnidadeMedida,
			&i.QuantidadePacote,
			&i.Peso,
			&i.ProdutoAtivo,
			&i.Ordem,
			&i.ProdutoImagem,
			&i.Categorias,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const pegarUsuarioEmail = `-- name: PegarUsuarioEmail :one
SELECT id, email, nome, telefone, role, ativo, created_at, updated_at
FROM usuarios
WHERE email = $1
LIMIT 1
`

func (q *Queries) PegarUsuarioEmail(ctx context.Context, email string) (Usuario, error) {
	row := q.db.QueryRow(ctx, pegarUsuarioEmail, email)
	var i Usuario
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Nome,
		&i.Telefone,
		&i.Role,
		&i.Ativo,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
