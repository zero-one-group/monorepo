import { type NextRequest, NextResponse } from 'next/server'

export function middleware(_req: NextRequest) {
  return NextResponse.next()
}

export const config = {
  matcher: ['/((?!api|_next/(?:data|static|image)|images|.*\\..*).*)'],
}
