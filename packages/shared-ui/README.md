# Shared UI Library

This is a shared UI library, a comprehensive collection of React components built with
modern web development best practices. Under the hood, it uses Tailwind CSS, Radix UI,
and TypeScript to deliver accessible, type-safe, and customizable components.

The UI components leverage [shadcn/ui](https://ui.shadcn.com/) as their foundation,
enhanced with custom modifications. The notable distinction lies in the styling
implementation, which utilizes [Tailwind Variants](https://www.tailwind-variants.org/)
for a more flexible approach and [tailwindcss-motion](https://rombo.co/tailwind) for
animations and transitions effects.

## Enabling the package

Add the following to your `package.json` file:

```json
{
    "dependencies": {
        "@repo/shared-ui": "workspace:*"
    },
    "devDependencies": {
        "@tailwindcss/vite": "^4.0.4",
        "tailwind-variants": "^0.3.1",
        "tailwindcss-motion": "^1.0.1",
        "tailwindcss": "^4.0.4"
    }
}
```

Add the following to your `tsconfig.json` file:

```json
{
    "references": [{ "path": "../../packages/shared-ui" }]
}
```

Exclude internal packages from optimization, add the following to your `vite.config.ts` file:

```ts
export default defineConfig({
  // ...
  optimizeDeps: {
    // Do not optimize internal workspace dependencies.
    exclude: ['@repo/shared-ui'],
  },
})
```

Finally, add the following to your `moon.yml` file:

```yaml
dependsOn:
  - 'shared-ui'
```
