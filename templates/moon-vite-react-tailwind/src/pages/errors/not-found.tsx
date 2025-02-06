export default function NotFound() {
  return (
    <main className="grid h-full min-h-screen place-items-center bg-white px-6 py-24 sm:py-32 lg:px-8">
      <div className="text-center">
        <p className="font-semibold text-base text-primary-600">404</p>
        <h1 className="mt-4 font-bold text-3xl text-neutral-900 tracking-tight sm:text-5xl">
          Page not found
        </h1>
        <p className="mt-6 text-base leading-7">
          Sorry, we couldn&apos;t find the page you&apos;re looking for.
        </p>
        <div className="mt-10 flex items-center justify-center gap-x-6">
          <a
            href="/"
            className="rounded-md bg-primary-600 px-3.5 py-2.5 font-semibold text-sm text-white shadow-sm hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-primary-600 focus-visible:outline-offset-2"
          >
            Go back home
          </a>
          <a
            href="https://example.com"
            className="font-semibold text-neutral-900 text-sm"
            target="_blank"
            rel="noreferrer"
          >
            Service Status <span aria-hidden="true">&rarr;</span>
          </a>
        </div>
      </div>
    </main>
  )
}
