import { execSync } from 'node:child_process'
import { existsSync, mkdirSync, unlinkSync, writeFileSync } from 'node:fs'
import { join } from 'node:path'
import {
  cancel,
  confirm,
  intro,
  isCancel,
  outro,
  select,
  spinner,
  tasks,
  text,
} from '@clack/prompts'
import { defineCommand } from 'citty'
import { createConsola } from 'consola'
import type { PhoenixAssetsOption, PhoenixDatabaseOption } from '../generator/phoenix'
import { generatePhoenixApp, promptPhoenixOptions } from '../generator/phoenix'

const _console = createConsola({ defaults: { tag: 'monorepo-cli' } })

// Define template info type for better type safety
type TemplateInfo = {
  template: string
  port: string
  requiresPort: boolean
}

type PhoenixInfo = {
  requiresPort: boolean
}

// Map of app types to their corresponding moon generate template commands and default ports
const TEMPLATE_MAP: Record<string, TemplateInfo | PhoenixInfo> = {
  astro: { template: 'template-astro', port: '4321', requiresPort: true },
  fastapi: { template: 'template-fastapi-ai', port: '5000', requiresPort: true },
  golang: { template: 'template-golang', port: '8080', requiresPort: true },
  nextjs: { template: 'template-nextjs', port: '3000', requiresPort: true },
  'react-spa': { template: 'template-react-app', port: '3000', requiresPort: true },
  'react-ssr': { template: 'template-react-ssr', port: '3000', requiresPort: true },
  'shared-ui': { template: 'template-shared-ui', port: '6300', requiresPort: false },
  strapi: { template: 'template-strapi', port: '1337', requiresPort: true },
  // Phoenix is not in the template map but we still define its port requirement
  phoenix: { requiresPort: false },
}

// Type guard to check if an entry has template and port
function isTemplateInfo(entry: TemplateInfo | PhoenixInfo): entry is TemplateInfo {
  return 'template' in entry && 'port' in entry
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
        initialValue: 'phoenix',
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
          { value: 'flutter', label: 'Flutter' },
        ],
      })

      if (appType === 'flutter') {
        const appName = await text({
          message: "What is the app's name?",
          placeholder: 'default : ZOG Mobile Starter',
          defaultValue: 'ZOG Mobile Starter',
        })
        if (isCancel(appName)) {
          cancel('Operation cancelled.')
          process.exit(0)
        }

        const appDirName = appName.toLowerCase().replace(/ /g, '-')
        const appDir = join(process.cwd(), 'apps', `flutter-${appDirName}`)
        if (!existsSync(appDir)) {
          mkdirSync(appDir, { recursive: true })
          _console.info(`Created folder: ${appDir}`)
        } else {
          _console.info(`Folder already exists: ${appDir}`)
        }

        const confirmAction = await confirm({
          message: `Do you want to generate the mobile flutter app in "${appDir}"?`,
          initialValue: true,
        })
        if (isCancel(confirmAction) || !confirmAction) {
          cancel('Operation cancelled.')
          process.exit(0)
        }

        const packageName = await text({
          message: "What is the app's package name?",
          placeholder: 'default : zog_starter',
          defaultValue: 'zog_starter',
        })
        if (isCancel(packageName)) {
          cancel('Operation cancelled.')
          process.exit(0)
        }

        const devFirebaseProjectId = await text({
          message: "What is the development app's Firebase project ID?",
          placeholder: 'default : zog-starter-dev',
          defaultValue: 'zog-starter-dev',
        })
        if (isCancel(devFirebaseProjectId)) {
          cancel('Operation cancelled.')
          process.exit(0)
        }

        const stgFirebaseProjectId = await text({
          message: "What is the staging app's Firebase project ID?",
          placeholder: 'default : zog-starter-stg',
          defaultValue: 'zog-starter-stg',
        })
        if (isCancel(stgFirebaseProjectId)) {
          cancel('Operation cancelled.')
          process.exit(0)
        }

        const prodFirebaseProjectId = await text({
          message: "What is the production app's Firebase project ID?",
          placeholder: 'default : zog-starter-prod',
          defaultValue: 'zog-starter-prod',
        })
        if (isCancel(prodFirebaseProjectId)) {
          cancel('Operation cancelled.')
          process.exit(0)
        }

        const androidAppId = await text({
          message: 'What is the Android App ID?',
          placeholder: 'default : com.zog.mobile',
          defaultValue: 'com.zog.mobile',
        })
        if (isCancel(androidAppId)) {
          cancel('Operation cancelled.')
          process.exit(0)
        }

        const iosBundleId = await text({
          message: 'What is the iOS Bundle ID?',
          placeholder: 'default : com.zog.mobile',
          defaultValue: 'com.zog.mobile',
        })
        if (isCancel(iosBundleId)) {
          cancel('Operation cancelled.')
          process.exit(0)
        }

        const s = spinner()
        s.start('Running mason commands...')

        await tasks([
          {
            title: 'Initialize mason',
            task: async () => {
              const masonYaml = join(appDir, 'mason.yaml')
              if (existsSync(masonYaml)) {
                unlinkSync(masonYaml)
              }
              execSync('mason init', { cwd: appDir, stdio: 'inherit' })
            },
          },
          {
            title: 'Add app_scaffolding brick',
            task: async () => {
              execSync(
                'mason add app_scaffolding --git-url https://github.com/zero-one-group/zog-mobile-bricks.git --git-path app_scaffolding',
                { cwd: appDir, stdio: 'inherit' }
              )
            },
          },
          {
            title: 'Run mason get',
            task: async () => {
              execSync('mason get', { cwd: appDir, stdio: 'inherit' })
            },
          },
          {
            title: 'Run mason make app_scaffolding',
            task: async () => {
              const config = {
                appName,
                packageName,
                devFirebaseProjectId,
                stgFirebaseProjectId,
                prodFirebaseProjectId,
                androidAppId,
                iosBundleId,
              }
              const configPath = join(appDir, 'config.json')
              writeFileSync(configPath, JSON.stringify(config, null, 2))

              try {
                execSync(
                  `mason make app_scaffolding --on-conflict overwrite --config-path config.json`,
                  { cwd: appDir, stdio: 'inherit' }
                )
              } finally {
                if (existsSync(configPath)) {
                  unlinkSync(configPath)
                }
              }
            },
          },
        ])

        s.stop('Mobile Flutter app generated!')
        outro(`You're all set! Your new mobile flutter app is ready in "${appDir}"! 🚀`)
        return
      }

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
          initialValue: true,
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

      // Get template info
      const templateInfo = TEMPLATE_MAP[appType] || { requiresPort: true }

      // Check if the app type requires a port
      const requiresPort = templateInfo.requiresPort

      // Get default port based on app type
      let defaultPort = args.portNumber
      if (isTemplateInfo(templateInfo)) {
        defaultPort = templateInfo.port
      }

      // Only ask for port number if the app type requires it
      let portNumber = defaultPort
      if (requiresPort) {
        const portResponse = await text({
          message: `What is the port number of your app? (ex: ${defaultPort})`,
          defaultValue: defaultPort,
        })

        if (isCancel(portResponse)) {
          cancel('Operation cancelled.')
          process.exit(0)
        }

        portNumber = portResponse
      }

      // Phoenix-specific options
      let phoenixOptions = {
        database: 'postgres' as PhoenixDatabaseOption,
        assets: 'both' as PhoenixAssetsOption,
        install: true,
      }

      if (appType === 'phoenix') {
        const options = await promptPhoenixOptions()
        if (!options) {
          process.exit(0)
        }
        phoenixOptions = options
      }

      // Build confirmation message
      let confirmMessage = `Do you want to generate the app with the following configuration?\n
  • Type: ${appType}
  • Name: ${appName}
  • Description: ${appDescription}\n`

      // Only include port in confirmation if it's required
      if (requiresPort) {
        confirmMessage += `\n  • Port: ${portNumber}`
      }

      // Add Phoenix-specific options to confirmation
      if (appType === 'phoenix') {
        confirmMessage += `\n  • Database: ${phoenixOptions.database}`
        confirmMessage += `\n  • Assets: ${phoenixOptions.assets}`
        confirmMessage += `\n  • Install dependencies: ${phoenixOptions.install ? 'Yes' : 'No'}`
      }

      // Add overwrite info if needed
      if (forceOverwrite) {
        confirmMessage += `\n  • Overwrite: Yes`
      }

      const confirmAction = await confirm({
        message: confirmMessage,
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
              _console.info('Generating Phoenix Framework application')
              try {
                await generatePhoenixApp({
                  appName,
                  appDescription,
                  database: phoenixOptions.database,
                  assets: phoenixOptions.assets,
                  force: forceOverwrite,
                  install: phoenixOptions.install,
                })
                return 'Phoenix Framework application generated successfully'
              } catch (error) {
                _console.error('Failed to generate Phoenix Framework application:', error)
                throw error
              }
            }

            const info = TEMPLATE_MAP[appType]
            if (info && isTemplateInfo(info)) {
              const templateName = info.template
              _console.info(`Generating ${appType} app using template: ${templateName}`)

              try {
                // Create command with arguments, adding --force flag if needed
                const forceFlag = forceOverwrite ? '--force' : ''
                let command = `moon generate ${templateName} ${forceFlag} -- \\
                  --package_name "${appName}" \\
                  --package_description "${appDescription}"`

                // Only add port number if required
                if (requiresPort) {
                  command += ` \\\n  --port_number "${portNumber}"`
                }

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
