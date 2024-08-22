# Shared UI Library

This is a shared UI library for the Monorepo project bsaed on `shadcn-ui`.

### Dependencies

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
  react-hook-form \
  react-resizable-panels \
  recharts \
  valibot \
  vaul
```

The key differences with the shadcn-ui library are the use of [`valibot`][valibot]
for validation library instead of zod, and the adoption of [`Modular Forms`][modular-form]
instead of `react-hook-form`.

<!-- link reference definition -->
[valibot]: https://valibot.dev/
[modular-form]: https://modularforms.dev/
