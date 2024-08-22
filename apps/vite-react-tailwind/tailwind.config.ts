import type { Config } from 'tailwindcss'
import colors from 'tailwindcss/colors'
import { fontFamily } from 'tailwindcss/defaultTheme'

export default {
  content: [
    './src/**/*!(*.stories|*.spec).{html,mdx,js,ts,jsx,tsx}',
    '../../packages/**/**!(*node_modules|*tests)/*!(*.stories|*.spec).{ts,tsx}',
    'index.html',
  ],
  darkMode: 'class',
  theme: {
    extend: {
      fontFamily: {
        sans: [...fontFamily.sans],
        mono: [...fontFamily.mono],
      },
      colors: {
        primary: colors.blue,
      },
    },
    debugScreens: {
      position: ['bottom', 'left'],
      borderTopRightRadius: '4px',
      printSize: true,
      prefix: '',
    },
  },
  plugins: [
    require('@tailwindcss/aspect-ratio'),
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
    require('tailwind-debug-breakpoints'),
    require('tailwindcss-animate'),
  ],
} satisfies Config
