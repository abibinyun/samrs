import { Hono } from 'hono'
import { cors } from 'hono/cors'

const app = new Hono()
const GO_SERVICE_URL = process.env.GO_SERVICE_URL || 'http://localhost:8080/api'

// Aktifkan CORS agar Vite bisa akses
app.use('*', cors({
  origin: ['http://localhost:5173', 'http://localhost:4173'],
  credentials: true,
}))

app.all('*', async (c) => {
  const url = new URL(c.req.url)
  
  // Gunakan template string sederhana untuk penggabungan yang lebih intuitif
  // url.pathname dari Vite adalah "/v1/auth/login"
  // GO_SERVICE_URL adalah "http://localhost:8080/api"
  const targetUrl = `${GO_SERVICE_URL}${url.pathname}${url.search}`

  console.log(`🚀 Bridge: ${c.req.method} ${url.pathname} -> ${targetUrl}`)

  const response = await fetch(targetUrl, {
    method: c.req.method,
    headers: c.req.header(),
    body: ['GET', 'HEAD'].includes(c.req.method) ? undefined : await c.req.arrayBuffer(),
  })

  // Salin response body dan headers secara manual untuk menghindari masalah immutability
  return new Response(response.body, {
    status: response.status,
    statusText: response.statusText,
    headers: response.headers
  })
})

export default app