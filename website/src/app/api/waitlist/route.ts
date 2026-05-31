import { NextResponse } from 'next/server'

export async function POST(req: Request) {
  try {
    const { name, email } = await req.json()
    if (!email || !email.includes('@')) {
      return NextResponse.json({ error: 'Valid email required' }, { status: 400 })
    }
    console.log('[Waitlist] New signup:', {
      name: name || 'anonymous',
      email,
      timestamp: new Date().toISOString(),
    })
    return NextResponse.json({ success: true })
  } catch {
    return NextResponse.json({ error: 'Invalid request' }, { status: 400 })
  }
}
