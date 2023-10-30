-- name: PegarUsuarioEmail :one
SELECT *
FROM usuarios
WHERE email = $1
LIMIT 1;
-- name: CriarUsuario :one
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
RETURNING *;
-- name: PegarTodosProdutos :many
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
LIMIT $1 OFFSET $2;
-- name: PegarProdutoUnico :one
SELECT p.id AS produto_id,
    p.nome AS produto_nome,
    p.descricao AS produto_descricao,
    p.preco AS produto_preco,
    p.unidade_medida AS produto_unidade_medida,
    p.quantidade_pacote AS produto_quantidade_pacote,
    p.peso AS produto_peso,
    p.ativo AS produto_ativo,
    p.ordem AS produto_ordem,
    p.imagem AS produto_imagem,
    array_to_string(array_agg(c.nome), ',') AS categorias_relacionadas
FROM produtos AS p
    JOIN "_CategoriaToProduto" AS cp ON p.id = cp."B"
    JOIN categorias AS c ON cp."A" = c.id
WHERE p.id = $1
GROUP BY p.id,
    p.nome,
    p.descricao,
    p.preco,
    p.unidade_medida,
    p.quantidade_pacote,
    p.peso,
    p.ativo,
    p.ordem,
    p.imagem;
-- name: CriarProduto :one
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
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, '')
RETURNING "id";
-- name: AtualizarProduto :exec
UPDATE "produtos"
SET "nome" = $1,
    "descricao" = $2,
    "preco" = $3,
    "unidade_medida" = $4,
    "quantidade_pacote" = $5,
    "peso" = $6,
    "ativo" = $7,
    "updated_at" = CURRENT_TIMESTAMP
WHERE "id" = $8;
-- name: AtualizarImagemProduto :exec
UPDATE "produtos"
SET "imagem" = 'https://firebasestorage.googleapis.com/v0/b/sos-do-maceneiro.appspot.com/o/produtos%2F' || "id" || '.webp?alt=media'
WHERE "id" = $1;
-- name: InserirCategoriaProdutoPeloNome :batchexec
INSERT INTO "_CategoriaToProduto" ("A", "B")
VALUES (
        (
            SELECT "id"
            FROM categorias
            WHERE 'nome' = $1
        ),
        (
            SELECT "id"
            FROM "produtos"
            WHERE 'id' = $2
        )
    ) ON CONFLICT ("A", "B") DO NOTHING;
-- name: PegarCategoriasProduto :many
SELECT c.*
FROM "categorias" c
    JOIN "_CategoriaToProduto" cp ON c."id" = cp."A"
WHERE cp."B" = $1;
-- name: PegarCategoriaPeloNome :batchmany
SELECT "id"
FROM "categorias"
WHERE "nome" = $1
LIMIT 1;
-- name: InserirCategoriaNoProduto :batchexec
INSERT INTO "_CategoriaToProduto" ("A", "B")
VALUES ($1, $2);
-- name: CriarCategoria :one
INSERT INTO "categorias" ("nome", "imagem")
VALUES ($1, $2)
RETURNING "id";
-- name: PegarTodasCategorias :many
SELECT *
FROM "categorias"
ORDER BY "updated_at" DESC
LIMIT $1 OFFSET $2;