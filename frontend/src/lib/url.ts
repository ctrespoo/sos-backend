
const url_base = import.meta.env.VITE_APP_URL ? '/api/v1' : 'http://localhost:20000/api/v1'

export const url = {
    produtos: `${url_base}/produtos`,
}
