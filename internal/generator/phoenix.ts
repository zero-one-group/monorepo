import { execSync } from 'node:child_process'
import { existsSync } from 'node:fs'
import { join } from 'node:path'
import { cancel, confirm, isCancel, select } from '@clack/prompts'
import { createConsola } from 'consola'

const _console = createConsola({ defaults: { tag: 'phoenix-generator' } })

/**
 * Database options for Phoenix Framework
 */
export type PhoenixDatabaseOption = 'postgres' | 'mysql' | 'mssql' | 'sqlite3' | 'none'

/**
 * Assets options for Phoenix Framework
 */
export type PhoenixAssetsOption = 'both' | 'esbuild' | 'tailwind' | 'none'

/**
 * Options for generating a Phoenix Framework application
 */
export interface PhoenixGenerateOptions {
  /**
   * The name of the OTP application
   */
  appName: string

  /**
   * The description of the application
   */
  appDescription: string

  /**
   * Specify the database adapter for Ecto
   * @default 'postgres'
   */
  database: PhoenixDatabaseOption

  /**
   * Specify assets to include
   * @default 'both'
   */
  assets: PhoenixAssetsOption

  /**
   * Force overwrite if the app directory already exists
   * @default false
   */
  force: boolean

  /**
   * Fetch and install dependencies
   * @default true
   */
  install: boolean
}

/**
 * Check if Erlang, Elixir, and Mix are installed
 * @returns true if all dependencies are installed, false otherwise
 */
function checkDependencies(): boolean {
  try {
    // Check Elixir version (which implicitly checks for Erlang)
    const elixirVersion = execSync('elixir --version', { encoding: 'utf8' })
    _console.info('Elixir detected:', elixirVersion.split('\n')[0])

    // Check Mix version
    const mixVersion = execSync('mix --version', { encoding: 'utf8' })
    _console.info('Mix detected:', mixVersion.trim())

    // Check if Phoenix generator is installed
    try {
      execSync('mix help phx.new', { stdio: 'ignore' })
      _console.info('Phoenix generator detected')
    } catch (_error) {
      _console.warn('Phoenix generator not found. Attempting to install...')
      execSync('mix archive.install hex phx_new', { stdio: 'inherit' })
      _console.success('Phoenix generator installed successfully')
    }

    return true
  } catch (_error) {
    _console.error('Required dependencies are missing:')
    _console.error('Please make sure Erlang and Elixir are installed.')
    _console.error('Installation instructions: https://elixir-lang.org/install.html')
    return false
  }
}

/**
 * Prompt for Phoenix Framework specific options
 * @returns Promise with Phoenix options or undefined if cancelled
 */
export async function promptPhoenixOptions(): Promise<
  | {
      database: PhoenixDatabaseOption
      assets: PhoenixAssetsOption
      install: boolean
    }
  | undefined
> {
  // Ask for database option
  const databaseOption = await select({
    message: 'Select database adapter for Ecto:',
    options: [
      { value: 'postgres', label: 'PostgreSQL' },
      { value: 'mysql', label: 'MySQL' },
      { value: 'mssql', label: 'MSSQL' },
      { value: 'sqlite3', label: 'SQLite3' },
      { value: 'none', label: 'No database (--no-ecto)' },
    ],
  })

  if (isCancel(databaseOption)) {
    cancel('Operation cancelled.')
    return undefined
  }

  // Ask for assets option
  const assetsOption = await select({
    message: 'Select assets to include:',
    options: [
      { value: 'both', label: 'Both esbuild and tailwind' },
      { value: 'esbuild', label: 'esbuild only (--no-tailwind)' },
      { value: 'tailwind', label: 'tailwind only (--no-esbuild)' },
      { value: 'none', label: 'No assets (--no-assets)' },
    ],
  })

  if (isCancel(assetsOption)) {
    cancel('Operation cancelled.')
    return undefined
  }

  // Ask if dependencies should be installed
  const installDeps = await confirm({
    message: 'Fetch and install dependencies?',
    initialValue: true,
  })

  if (isCancel(installDeps)) {
    cancel('Operation cancelled.')
    return undefined
  }

  return {
    database: databaseOption as PhoenixDatabaseOption,
    assets: assetsOption as PhoenixAssetsOption,
    install: installDeps,
  }
}

/**
 * Generate a Phoenix Framework application
 * @param options Options for generating a Phoenix Framework application
 * @returns Promise that resolves when the application is generated
 */
export async function generatePhoenixApp(options: PhoenixGenerateOptions): Promise<void> {
  const { appName, database, assets, force, install } = options

  _console.info(`Generating Phoenix Framework application: ${appName}`)

  // Check if required dependencies are installed
  if (!checkDependencies()) {
    throw new Error('Required dependencies for Phoenix Framework are not installed')
  }

  // Create apps directory if it doesn't exist
  const appsDir = join(process.cwd(), 'apps')
  execSync(`mkdir -p ${appsDir}`, { stdio: 'inherit' })

  // Handle force option - remove existing directory if force is true
  const appDir = join(appsDir, appName)
  if (existsSync(appDir) && force) {
    _console.warn(`Removing existing directory: ${appDir}`)
    execSync(`rm -rf ${appDir}`, { stdio: 'inherit' })
  }

  // Build command arguments
  const args: string[] = []

  // Add app name
  args.push(appName)

  // Add database option
  if (database === 'none') {
    args.push('--no-ecto')
  } else {
    args.push(`--database ${database}`)
  }

  // Add assets options
  if (assets === 'none') {
    args.push('--no-assets')
  } else if (assets === 'esbuild') {
    args.push('--no-tailwind')
  } else if (assets === 'tailwind') {
    args.push('--no-esbuild')
  }

  // Add install option - only add --no-install if install is false
  if (!install) {
    args.push('--no-install')
  } else {
    args.push('--install')
  }

  // Build the full command
  const command = `mix phx.new ${args.join(' ')}`

  try {
    // Change to apps directory and execute the command
    process.chdir(appsDir)
    _console.info(`Executing command: ${command}`)
    execSync(command, { stdio: 'inherit' })

    _console.success(`Phoenix Framework application ${appName} generated successfully!`)

    if (!install) {
      _console.success(`Run \`moon ${appName}:deps\` to install dependencies.`)
    }
  } catch (error) {
    _console.error('Failed to generate Phoenix Framework application:', error)
    throw error
  }
}
