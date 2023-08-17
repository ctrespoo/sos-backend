package model

type RespErro struct {
	Erro     int    `json:"erro,omitempty"`
	Mensagem string `json:"mensagem"`
}

type RespSucesso struct {
	Token string `json:"token"`
}
