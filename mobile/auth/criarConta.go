package mobile

import (
	"encoding/json"
	"log"
	"net/http"

	model "sos/backend/model"

	interno "sos/backend/interno/db"

	"firebase.google.com/go/v4/auth"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Nome        string `json:"nome"`
	Email       string `json:"email"`
	Telefone    string `json:"telefone"`
	Senha       string `json:"senha"`
	RepitaSenha string `json:"repitaSenha"`
}

func CriarConta(db *interno.Queries, client *auth.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		var user user
		json.NewDecoder(r.Body).Decode(&user)

		if user.Nome == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&model.RespErro{Erro: 400, Mensagem: "Nome não pode ser vazia"})
			return
		}
		if user.Telefone == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&model.RespErro{Erro: 400, Mensagem: "Telefone não pode ser vazio"})
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
		if len(user.Senha) < 6 {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&model.RespErro{Erro: 400, Mensagem: "Senha precisa ter no minimo 6 caracteres"})
			return
		}

		// Cria usuario no postgres
		hashSenha, err := bcrypt.GenerateFromPassword([]byte(user.Senha), 14)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&model.RespErro{Erro: 500, Mensagem: "Erro ao gerar hash da senha"})
			return
		}

		res, err := db.CriarUsuario(r.Context(), interno.CriarUsuarioParams{
			Telefone: user.Telefone,
			Nome:     user.Nome,
			Email:    user.Email,
			Senha:    string(hashSenha),
		})
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(&model.RespErro{Erro: 500, Mensagem: "Erro ao criar usuario postgres"})
			return
		}

		// Cria usuario no firebase
		params := (&auth.UserToCreate{}).
			PhoneNumber(user.Telefone).
			UID(res.ID.String()).
			Email(res.Email).
			EmailVerified(false).
			Password(user.Senha).
			DisplayName(res.Nome).
			Disabled(false)

		u, err := client.CreateUser(r.Context(), params)
		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&model.RespErro{Erro: 500, Mensagem: err.Error()})
			return
		}

		claims := map[string]interface{}{
			"teste": "testando",
		}

		err = client.SetCustomUserClaims(r.Context(), u.UID, claims)
		if err != nil {
			log.Println(err)
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&model.RespErro{Erro: 500, Mensagem: "Erro token"})
			return
		}

		w.WriteHeader(201)
		w.Write([]byte("Sucesso"))
	}
}
