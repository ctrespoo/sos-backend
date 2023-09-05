-- Crie 10 produtos aleatórios
WITH random_products AS (
    SELECT 'Produto ' || generate_series(1, 10) AS nome,
        'Descrição do Produto ' || generate_series(1, 10) AS descricao,
        (random() * 100)::numeric(15, 2) AS preco,
        CASE
            WHEN random() < 0.5 THEN 'KG'
            ELSE 'G'
        END AS unidade_medida,
        (random() * 10)::integer AS quantidade_pacote,
        (random() * 1000)::numeric(15, 2) AS peso,
        CASE
            WHEN random() < 0.5 THEN true
            ELSE false
        END AS ativo,
        (random() * 100)::integer AS ordem,
        'Imagem do Produto ' || generate_series(1, 10) AS imagem
)
INSERT INTO produtos (
        nome,
        descricao,
        preco,
        unidade_medida,
        quantidade_pacote,
        peso,
        ativo,
        ordem,
        imagem,
        created_at,
        updated_at
    )
SELECT nome,
    descricao,
    preco,
    unidade_medida::"TipoUnidadeMedida",
    quantidade_pacote,
    peso,
    ativo,
    ordem,
    imagem,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
FROM random_products;
-- Crie categorias aleatórias para cada produto
INSERT INTO categorias (nome, imagem, created_at, updated_at)
SELECT 'Categoria ' || generate_series(1, 10) AS nome,
    'Imagem da Categoria ' || generate_series(1, 10) AS imagem,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP;
-- Vincule produtos a categorias aleatórias
INSERT INTO "_CategoriaToProduto" ("A", "B")
SELECT c.id AS categoria_id,
    p.id AS produto_id
FROM categorias c
    JOIN (
        SELECT id
        FROM produtos
        ORDER BY random()
        LIMIT 10
    ) p ON true;