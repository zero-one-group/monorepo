import type { NextRequest } from 'next/server'

export async function GET(_req: NextRequest) {
  const headers: HeadersInit = { 'Content-Type': 'application/rss+xml' }
  return new Response('rss.xml', { status: 200, headers })
}
