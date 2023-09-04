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
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING "id";
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