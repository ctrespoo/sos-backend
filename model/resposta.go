package model

type RespErro struct {
	Erro     int    `json:"erro,omitempty"`
	Mensagem string `json:"mensagem"`
}
