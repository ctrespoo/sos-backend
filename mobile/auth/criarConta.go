package mobile

import (
	"encoding/json"
	"log"
	"net/http"

	model "sos/backend/model"

	interno "sos/backend/interno/db"

	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Nome        string `json:"nome"`
	Email       string `json:"email"`
	Senha       string `json:"senha"`
	RepitaSenha string `json:"repitaSenha"`
}

func CriarConta(db *interno.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		var user user
		json.NewDecoder(r.Body).Decode(&user)

		if user.Nome == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&model.RespErro{Erro: 400, Mensagem: "Nome não pode ser vazia"})
			return
		}
		if user.Email == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&model.RespErro{Erro: 400, Mensagem: "Usuario não pode ser vazia"})
			return
		}
		if user.Senha == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&model.RespErro{Erro: 400, Mensagem: "Senha não pode ser vazia"})
			return
		}
		if user.Senha != user.RepitaSenha {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&model.RespErro{Erro: 400, Mensagem: "Senhas não conferem"})
			return
		}

		hashSenha, err := bcrypt.GenerateFromPassword([]byte(user.Senha), 14)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&model.RespErro{Erro: 500, Mensagem: "Erro ao gerar hash da senha"})
			return
		}

		res, err := db.CriarUsuario(r.Context(), interno.CriarUsuarioParams{
			Nome:  user.Nome,
			Email: user.Email,
			Senha: string(hashSenha),
		})
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&model.RespErro{Erro: 500, Mensagem: "Erro ao criar usuario"})
			return
		}
		log.Println(res)
		w.Write([]byte("criarConta!"))
	}
}
