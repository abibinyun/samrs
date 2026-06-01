import { Hono } from 'hono'
import { cors } from 'hono/cors'

const app = new Hono()
const GO_SERVICE_URL = process.env.GO_SERVICE_URL || 'http://localhost:8080'

app.use('*', cors({
  origin: ['http://localhost:5173', 'http://localhost:4173'],
  credentials: true,
}))

app.all('*', async (c) => {
  const url = new URL(c.req.url)
  const pathname = url.pathname

  // Forward path as-is to backend
  // Dev (Vite proxy strips /api): receives /v1/auth/login → backend expects /api/v1/auth/login
  // Prod (nginx preserves /api): receives /api/v1/auth/login → backend expects /api/v1/auth/login
  const apiPath = pathname.startsWith('/api') ? pathname : `/api${pathname}`
  const targetUrl = `${GO_SERVICE_URL}${apiPath}${url.search}`

  const response = await fetch(targetUrl, {
    method: c.req.method,
    headers: c.req.header(),
    body: ['GET', 'HEAD'].includes(c.req.method) ? undefined : await c.req.arrayBuffer(),
  })

  return new Response(response.body, {
    status: response.status,
    statusText: response.statusText,
    headers: response.headers
  })
})

export default app