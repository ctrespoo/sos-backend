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