import { spawn } from 'node:child_process'
import { intro, outro, spinner, tasks } from '@clack/prompts'
import { cancel, confirm, isCancel, select, text } from '@clack/prompts'
import { defineCommand } from 'citty'
import { createConsola } from 'consola'

const _console = createConsola({ defaults: { tag: 'monorepo-cli' } })

export default defineCommand({
  meta: {
    name: 'generate:app',
    description: 'Generate new application',
  },
  args: {
    appName: {
      type: 'string',
      description: 'The name of the app',
      default: 'my-app',
    },
    appDescription: {
      type: 'string',
      description: 'The description of the app',
      default: 'My awesome application',
    },
    portNumber: {
      type: 'string',
      description: 'The port number of the app',
      default: '3000',
    },
    help: {
      type: 'boolean',
      description: 'Print information about the command',
      default: false,
    },
  },
  async run({ args }) {
    try {
      intro(`Creating new application...`)

      const appType = await select({
        message: 'Pick a project type.',
        options: [
          { value: 'astro', label: 'Astro' },
          { value: 'fastapi', label: 'FastAPI' },
          { value: 'golang', label: 'Golang' },
          { value: 'nextjs', label: 'Next.js' },
          { value: 'phoenix', label: 'Phoenix Framework' },
          { value: 'react-spa', label: 'React SPA' },
          { value: 'react-ssr', label: 'React SSR' },
          { value: 'shared-ui', label: 'Shared UI' },
          { value: 'strapi', label: 'Strapi' },
        ],
      })

      if (isCancel(appType)) {
        cancel('Operation cancelled.')
        process.exit(0)
      }

      // Use snake_case for Phoenix apps
      const defaultAppName = appType === 'phoenix' ? args.appName.replace('-', '_') : args.appName

      const appName = await text({
        message: `What is the name of your app? (ex: ${defaultAppName})`,
        defaultValue: defaultAppName,
      })

      if (isCancel(appName)) {
        cancel('Operation cancelled.')
        process.exit(0)
      }

      const appTitle = await text({
        message: `What is the title of your app? (ex: My Application)`,
        defaultValue: 'My Application',
      })

      if (isCancel(appTitle)) {
        cancel('Operation cancelled.')
        process.exit(0)
      }

      const appDescription = await text({
        message: `What is the description of your app?`,
        defaultValue: args.appDescription,
      })

      if (isCancel(appDescription)) {
        cancel('Operation cancelled.')
        process.exit(0)
      }

      const portNumber = await text({
        message: `What is the port number of your app? (ex: ${args.portNumber})`,
        defaultValue: args.portNumber,
      })

      if (isCancel(portNumber)) {
        cancel('Operation cancelled.')
        process.exit(0)
      }

      const confirmAction = await confirm({
        message: `Do you want to generate the app with the following configuration?`,
        initialValue: true,
      })

      if (isCancel(confirmAction)) {
        cancel('Operation cancelled.')
        process.exit(0)
      }

      // Build the params object
      const params = { appName, appTitle, appDescription, portNumber }

      const s = spinner()
      s.start('Generating app...')

      _console.info(params)

      await tasks([
        {
          title: 'Generating application files',
          task: async (message) => {
            // Do installation here

            // moon setup
            // moon generate template-strapi

            const ls = spawn('ls', ['-lh', '/Users/ariss/Developer/Excercise'])

            ls.stdout.on('data', (data) => {
              _console.log(`stdout: ${data}`)
            })

            ls.stderr.on('data', (data) => {
              _console.error(`stderr: ${data}`)
            })

            ls.on('close', (code) => {
              _console.log(`child process exited with code ${code}`)
            })
            return 'Installed via npm'
          },
        },
      ])

      s.stop('New application generated!')
      outro(`You're all set!`)
    } catch (error) {
      _console.error(error instanceof Error ? error.message : 'Unknown error occurred')
      process.exit(1)
    }
  },
})
