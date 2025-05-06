import 'dotenv/config'
import { defineCommand, runMain, showUsage } from 'citty'
import pkg from '~~/package.json' assert { type: 'json' }

const main = defineCommand({
  meta: {
    name: 'cmd',
    version: pkg.version,
    description: `Monorepo Command Line Interface`,
  },
  args: {
    help: {
      type: 'boolean',
      description: 'Print information about the application',
      default: false,
    },
  },
  subCommands: {
    'generate:app': () => import('./cmd/generate-app').then((r) => r.default),
    'make:app-key': () => import('./cmd/make-app-key').then((r) => r.default),
    'make:strapi-keys': () => import('./cmd/make-strapi-keys').then((r) => r.default),
  },
  async run({ args, cmd }) {
    // Show help page if --help flag is used or no subcommand provided
    if (args.help || args._.length === 0) {
      showUsage(cmd)
      return
    }
  },
})

runMain(main)
