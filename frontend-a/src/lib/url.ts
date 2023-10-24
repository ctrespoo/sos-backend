
const url_base = import.meta.env.APP_URL ? '/api/v1' : 'http://localhost:20000/api/v1'

export const url = {
    criarConta: `${url_base}/auth/registrar`,
    produtos: `${url_base}/produtos`,
}