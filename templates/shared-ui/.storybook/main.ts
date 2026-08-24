import type { StorybookConfig } from '@storybook/react-vite'
import { mergeConfig } from 'vite'

const config: StorybookConfig = {
  stories: ['./_docs/**/*.mdx', '../src/**/*.mdx', '../src/**/*.stories.@(ts|tsx)'],

  addons: ['@storybook/addon-links', '@storybook/addon-a11y', '@storybook/addon-docs'],

  framework: {
    name: '@storybook/react-vite',
    options: {},
  },

  core: {
    disableTelemetry: true, // 👈 Disables telemetry
    enableCrashReports: false, // 👈 Appends the crash reports to the telemetry events
    disableWhatsNewNotifications: true, // 👈 Disables the whats new notification
  },

  async viteFinal(config) {
    return mergeConfig(config, {
      plugins: [],
      build: {
        chunkSizeWarningLimit: 1024 * 4,
      },
    })
  },

  features: {
    backgrounds: false,
  },
}

export default config
