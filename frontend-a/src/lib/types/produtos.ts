export interface Produto {
    ProdutoID: string
    ProdutoNome: string
    ProdutoDescricao: string
    Preco: number
    UnidadeMedida: string
    QuantidadePacote: number
    Peso: number
    ProdutoAtivo: boolean
    Ordem: number
    ProdutoImagem: string
    Categorias: string
}

export interface Categoria {
    ID: string
    Nome: string
    Imagem: string
    CreatedAt: string
    UpdatedAt: string
}

export interface ProdutoUnico {
    ProdutoID: string,
    ProdutoNome: string,
    ProdutoDescricao: string,
    ProdutoPreco: number,
    ProdutoUnidadeMedida: string,
    ProdutoQuantidadePacote: number,
    ProdutoPeso: number,
    ProdutoAtivo: boolean,
    ProdutoOrdem: number,
    ProdutoImagem: string,
    CategoriasRelacionadas: string
}