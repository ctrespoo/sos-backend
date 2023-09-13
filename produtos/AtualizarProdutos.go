package produtos

import (
	"net/http"
	interno "sos/backend/interno/db"

	"firebase.google.com/go/v4/auth"
)

func AtualizarProduto(db *interno.Queries, app *auth.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var categorias []interno.Categoria
		var nomesCategorias []string
		var categoriasParaExcluir []string

		verificarCategorias(&categorias, &nomesCategorias, &categoriasParaExcluir)

	}
}

func verificarCategorias(categorias *[]interno.Categoria, nomesCategorias *[]string, categoriasParaExcluir *[]string) {
	// Percorra as categorias
	for _, categoria := range *categorias {
		// Verifique se o nome da categoria não está na lista de nomesCategorias
		if !contemNomeCategoria(categoria.Nome, nomesCategorias) {
			// Adicione o nome da categoria à lista de categoriasParaExcluir
			*categoriasParaExcluir = append(*categoriasParaExcluir, categoria.Nome)
		}
	}
}

func contemNomeCategoria(nome string, nomesCategorias *[]string) bool {
	// Percorra a lista de nomesCategorias e verifique se o nome está presente
	for _, nomeCategoria := range *nomesCategorias {
		if nomeCategoria == nome {
			return true
		}
	}
	return false
}
