import { execSync } from 'node:child_process'
import { existsSync } from 'node:fs'
import { join } from 'node:path'
import { intro, outro, spinner, tasks } from '@clack/prompts'
import { cancel, confirm, isCancel, select, text } from '@clack/prompts'
import { defineCommand } from 'citty'
import { createConsola } from 'consola'

const _console = createConsola({ defaults: { tag: 'monorepo-cli' } })

// Map of app types to their corresponding moon generate template commands
const TEMPLATE_MAP = {
  astro: 'template-astro',
  fastapi: 'template-fastapi-ai',
  golang: 'template-golang',
  nextjs: 'template-nextjs',
  'react-spa': 'template-react-app',
  'react-ssr': 'template-react-ssr',
  'shared-ui': 'template-shared-ui',
  strapi: 'template-strapi',
}

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

      // Check if app directory already exists
      const appDir = join(process.cwd(), 'apps', appName)
      const appExists = existsSync(appDir)
      let forceOverwrite = false

      if (appExists) {
        const overwriteConfirm = await confirm({
          message: `App directory "${appName}" already exists. Do you want to overwrite it?`,
          initialValue: false,
        })

        if (isCancel(overwriteConfirm) || !overwriteConfirm) {
          cancel('Operation cancelled.')
          process.exit(0)
        }

        forceOverwrite = true
      }

      const appDescription = await text({
        message: `What is the description of your app? (ex: ${args.appDescription})`,
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
        message: `Do you want to generate the app with the following configuration?
  • Type: ${appType}
  • Name: ${appName}
  • Description: ${appDescription}
  • Port: ${portNumber}${forceOverwrite ? '\n  • Overwrite: Yes' : ''}`,
        initialValue: true,
      })

      if (isCancel(confirmAction) || !confirmAction) {
        cancel('Operation cancelled.')
        process.exit(0)
      }

      const s = spinner()
      s.start('Generating app...')

      await tasks([
        {
          title: 'Generating application files',
          task: async () => {
            // Generate app based on selected type
            if (appType === 'phoenix') {
              // Phoenix has a special case since it's not in the template map
              _console.info('Phoenix generator is not implemented yet')
              throw new Error('Phoenix generator is not implemented yet')
            }

            if (appType in TEMPLATE_MAP) {
              const templateName = TEMPLATE_MAP[appType as keyof typeof TEMPLATE_MAP]
              _console.info(`Generating ${appType} app using template: ${templateName}`)

              try {
                // Create command with arguments, adding --force flag if needed
                const forceFlag = forceOverwrite ? '--force' : ''
                const command = `moon generate ${templateName} ${forceFlag} -- \\
                  --package_name "${appName}" \\
                  --package_description "${appDescription}" \\
                  --port_number "${portNumber}"`

                // Execute the command
                execSync(command, { stdio: 'inherit' })

                return `${appType} application generated successfully`
              } catch (error) {
                _console.error(`Failed to generate ${appType} application:`, error)
                throw error
              }
            }

            throw new Error(`Generator for ${appType} is not implemented yet`)
          },
        },
      ])

      s.stop('New application generated!')
      outro(`You're all set! Your new ${appType} application "${appName}" is ready.`)
    } catch (error) {
      _console.error(error instanceof Error ? error.message : 'Unknown error occurred')
      process.exit(1)
    }
  },
})
