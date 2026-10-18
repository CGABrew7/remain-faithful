import { proxyAdmin } from '@/lib/admin-proxy'

const ALLOWED = new Set([
  'login',
  'logout',
  'overview',
  'health',
  'account',
  'account/resend-verification',
  'account/password-reset',
  'account/delete',
])

async function handle(request: Request, ctx: { params: Promise<{ path: string[] }> }) {
  const { path } = await ctx.params
  const joined = path.join('/')
  if (!ALLOWED.has(joined)) {
    return Response.json({ error: 'not found' }, { status: 404 })
  }
  return proxyAdmin(request, joined)
}

export const GET = handle
export const POST = handle
