'use client'

export default function Template({ children }: React.PropsWithChildren) {
  return <div className="flex min-h-screen flex-col pt-16 pb-12">{children}</div>
}
