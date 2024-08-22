# Shared UI Library

This is a shared UI library for the Monorepo project based on `shadcn-ui`.

## Dependencies

### RadixUI Packages

```sh
pnpm --filter @myorg/shared-ui add \
  @radix-ui/react-accordion \
  @radix-ui/react-alert-dialog \
  @radix-ui/react-aspect-ratio \
  @radix-ui/react-avatar \
  @radix-ui/react-checkbox \
  @radix-ui/react-collapsible \
  @radix-ui/react-context-menu \
  @radix-ui/react-dialog \
  @radix-ui/react-dropdown-menu \
  @radix-ui/react-hover-card \
  @radix-ui/react-icons \
  @radix-ui/react-label \
  @radix-ui/react-menubar \
  @radix-ui/react-navigation-menu \
  @radix-ui/react-popover \
  @radix-ui/react-progress \
  @radix-ui/react-radio-group \
  @radix-ui/react-scroll-area \
  @radix-ui/react-select \
  @radix-ui/react-separator \
  @radix-ui/react-slider \
  @radix-ui/react-slot \
  @radix-ui/react-switch \
  @radix-ui/react-tabs \
  @radix-ui/react-toast \
  @radix-ui/react-toggle \
  @radix-ui/react-toggle-group \
  @radix-ui/react-tooltip
```

### Additional Packages

```sh
pnpm --filter @myorg/shared-ui add \
  cmdk \
  embla-carousel-react \
  input-otp \
  react-day-picker \
  react-resizable-panels \
  recharts \
  vaul
```

The difference with `shadcn-ui` is that it does not include the `Form` component
from `react-hook-form`, the `Data Table` component that uses `TanStack Table`, and
the `zod` validation library. This is intended to make it more flexible to use
different form and validation libraries, for example if you want to use
[`Modular Forms`][modular-form] and [`valibot`][valibot].


## Tailwind CSS Config

Add the following to your `tailwind.config.js` file:

```ts
import type { Config } from 'tailwindcss'
import { fontFamily } from 'tailwindcss/defaultTheme'

export default {
  // ... other tailwind config
  theme: {
    extend: {
      fontFamily: {
        sans: [...fontFamily.sans],
        mono: [...fontFamily.mono],
      },
    },
    container: {
      center: true,
      padding: '2rem',
      screens: {
        '2xl': '1400px',
      },
    },
    colors: {
      border: 'hsl(var(--border))',
      input: 'hsl(var(--input))',
      ring: 'hsl(var(--ring))',
      background: 'hsl(var(--background))',
      foreground: 'hsl(var(--foreground))',
      primary: {
        DEFAULT: 'hsl(var(--primary))',
        foreground: 'hsl(var(--primary-foreground))',
      },
      secondary: {
        DEFAULT: 'hsl(var(--secondary))',
        foreground: 'hsl(var(--secondary-foreground))',
      },
      destructive: {
        DEFAULT: 'hsl(var(--destructive))',
        foreground: 'hsl(var(--destructive-foreground))',
      },
      muted: {
        DEFAULT: 'hsl(var(--muted))',
        foreground: 'hsl(var(--muted-foreground))',
      },
      accent: {
        DEFAULT: 'hsl(var(--accent))',
        foreground: 'hsl(var(--accent-foreground))',
      },
      popover: {
        DEFAULT: 'hsl(var(--popover))',
        foreground: 'hsl(var(--popover-foreground))',
      },
      card: {
        DEFAULT: 'hsl(var(--card))',
        foreground: 'hsl(var(--card-foreground))',
      },
    },
    borderRadius: {
      lg: `var(--radius)`,
      md: `calc(var(--radius) - 2px)`,
      sm: 'calc(var(--radius) - 4px)',
    },
    keyframes: {
      'accordion-down': {
        from: { height: '0' },
        to: { height: 'var(--radix-accordion-content-height)' },
      },
      'accordion-up': {
        from: { height: 'var(--radix-accordion-content-height)' },
        to: { height: '0' },
      },
      'caret-blink': {
        '0%,70%,100%': { opacity: '1' },
        '20%,50%': { opacity: '0' },
      },
    },
    animation: {
      'accordion-down': 'accordion-down 0.2s ease-out',
      'accordion-up': 'accordion-up 0.2s ease-out',
      'caret-blink': 'caret-blink 1.25s ease-out infinite',
    },
    // ...  other theme configuration
  },
  plugins: [
    require('tailwindcss-animate'),
    // ... other plugins
  ],
} satisfies Config
```

And add the following styles to your `globals.css` file:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  :root {
    --background: 0 0% 100%;
    --foreground: 240 10% 3.9%;
    --card: 0 0% 100%;
    --card-foreground: 240 10% 3.9%;
    --popover: 0 0% 100%;
    --popover-foreground: 240 10% 3.9%;
    --primary: 240 5.9% 10%;
    --primary-foreground: 0 0% 98%;
    --secondary: 240 4.8% 95.9%;
    --secondary-foreground: 240 5.9% 10%;
    --muted: 240 4.8% 95.9%;
    --muted-foreground: 240 3.8% 46.1%;
    --accent: 240 4.8% 95.9%;
    --accent-foreground: 240 5.9% 10%;
    --destructive: 0 72.22% 50.59%;
    --destructive-foreground: 0 0% 98%;
    --border: 240 5.9% 90%;
    --input: 240 5.9% 90%;
    --ring: 240 5% 64.9%;
    --radius: 0.5rem;

    --chart-1: 12 76% 61%;
    --chart-2: 173 58% 39%;
    --chart-3: 197 37% 24%;
    --chart-4: 43 74% 66%;
    --chart-5: 27 87% 67%;
  }

  .dark {
    --background: 240 10% 3.9%;
    --foreground: 0 0% 98%;
    --card: 240 10% 3.9%;
    --card-foreground: 0 0% 98%;
    --popover: 240 10% 3.9%;
    --popover-foreground: 0 0% 98%;
    --primary: 0 0% 98%;
    --primary-foreground: 240 5.9% 10%;
    --secondary: 240 3.7% 15.9%;
    --secondary-foreground: 0 0% 98%;
    --muted: 240 3.7% 15.9%;
    --muted-foreground: 240 5% 64.9%;
    --accent: 240 3.7% 15.9%;
    --accent-foreground: 0 0% 98%;
    --destructive: 0 62.8% 30.6%;
    --destructive-foreground: 0 85.7% 97.3%;
    --border: 240 3.7% 15.9%;
    --input: 240 3.7% 15.9%;
    --ring: 240 4.9% 83.9%;

    --chart-1: 220 70% 50%;
    --chart-2: 160 60% 45%;
    --chart-3: 30 80% 55%;
    --chart-4: 280 65% 60%;
    --chart-5: 340 75% 55%;
  }
}

@layer base {
  * {
    @apply border-border;
  }
  html,
  body {
    @apply size-full min-h-screen antialiased scroll-smooth bg-background text-foreground;
  }
  body {
    font-feature-settings: "rlig" 1, "calt" 1;
  }
}
```

## Other Configurations

Add the following to your `vite.config.ts` file:

```ts
export default defineConfig({
  // ....
  optimizeDeps: {
    // Do not optimize internal workspace dependencies.
    exclude: ['@myorg/shared-ui'],
  },
  resolve: {
    alias: [
      {
        // Configure an alias for the package. So, we don't have to restart
        // the Vite server each time when the former is performed.
        find: '@myorg/shared-ui',
        replacement: resolve(__dirname, '../../packages/shared-ui/src/index.ts'),
      },
    ],
  },
  // ....
})
```

Add the following to your `package.json` file:

```json
{
    "dependencies": {
        "@myorg/shared-ui": "workspace:*",
        "tailwindcss-animate": "^1.0.7",
    }
}
```

Add the following to your `tsconfig.json` file:

```json
{
    "references": [{ "path": "../../packages/shared-ui" } ]
}
```

Finally, add the following to your `moon.yml` file:

```yaml
dependsOn:
  - 'shared-ui'
```

<!-- link reference definition -->
[valibot]: https://valibot.dev/
[modular-form]: https://modularforms.dev/
