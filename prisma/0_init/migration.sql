-- CreateEnum
CREATE TYPE "TipoUnidadeMedida" AS ENUM ('KG', 'G');

-- CreateTable
CREATE TABLE "usuarios" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "email" TEXT NOT NULL,
    "nome" TEXT NOT NULL,
    "telefone" TEXT NOT NULL,
    "senha" TEXT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "usuarios_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "categorias" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "nome" TEXT NOT NULL,
    "imagem" TEXT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "categorias_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "loja" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),

    CONSTRAINT "loja_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "produtos" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "nome" TEXT NOT NULL,
    "descricao" TEXT NOT NULL,
    "preco" DECIMAL(15,2) NOT NULL,
    "unidade_medida" "TipoUnidadeMedida" NOT NULL,
    "quantidade_pacote" INTEGER NOT NULL,
    "peso" DECIMAL(15,2) NOT NULL,
    "ativo" BOOLEAN NOT NULL DEFAULT false,
    "ordem" INTEGER NOT NULL DEFAULT 0,
    "imagem" TEXT NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "produtos_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "produtos_loja" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "sku" TEXT,
    "barcode" TEXT,
    "produto_id" UUID NOT NULL,
    "loja_id" UUID NOT NULL,
    "ativo" BOOLEAN DEFAULT false,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "produtos_loja_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "_CategoriaToProduto" (
    "A" UUID NOT NULL,
    "B" UUID NOT NULL
);

-- CreateIndex
CREATE UNIQUE INDEX "usuarios_email_key" ON "usuarios"("email");

-- CreateIndex
CREATE UNIQUE INDEX "usuarios_telefone_key" ON "usuarios"("telefone");

-- CreateIndex
CREATE UNIQUE INDEX "categorias_nome_key" ON "categorias"("nome");

-- CreateIndex
CREATE UNIQUE INDEX "produtos_loja_produto_id_loja_id_key" ON "produtos_loja"("produto_id", "loja_id");

-- CreateIndex
CREATE UNIQUE INDEX "produtos_loja_loja_id_sku_key" ON "produtos_loja"("loja_id", "sku");

-- CreateIndex
CREATE UNIQUE INDEX "_CategoriaToProduto_AB_unique" ON "_CategoriaToProduto"("A", "B");

-- CreateIndex
CREATE INDEX "_CategoriaToProduto_B_index" ON "_CategoriaToProduto"("B");

-- AddForeignKey
ALTER TABLE "produtos_loja" ADD CONSTRAINT "produtos_loja_produto_id_fkey" FOREIGN KEY ("produto_id") REFERENCES "produtos"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "produtos_loja" ADD CONSTRAINT "produtos_loja_loja_id_fkey" FOREIGN KEY ("loja_id") REFERENCES "loja"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_CategoriaToProduto" ADD CONSTRAINT "_CategoriaToProduto_A_fkey" FOREIGN KEY ("A") REFERENCES "categorias"("id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "_CategoriaToProduto" ADD CONSTRAINT "_CategoriaToProduto_B_fkey" FOREIGN KEY ("B") REFERENCES "produtos"("id") ON DELETE CASCADE ON UPDATE CASCADE;

