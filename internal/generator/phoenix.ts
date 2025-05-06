import { execSync } from 'node:child_process'
import { existsSync, readFileSync, writeFileSync } from 'node:fs'
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
      { value: 'none', label: 'No database (--no-ecto)' },
      { value: 'postgres', label: 'PostgreSQL' },
      { value: 'mysql', label: 'MySQL' },
      { value: 'mssql', label: 'MSSQL' },
      { value: 'sqlite3', label: 'SQLite3' },
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
    initialValue: false,
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
 * Ensures dependencies are installed before running commands that require them
 * @param appDir Directory of the Phoenix application
 * @param userInstallChoice Whether the user chose to install dependencies
 * @returns true if dependencies are installed or were installed successfully
 */
function ensureDependenciesInstalled(appDir: string, userInstallChoice: boolean): boolean {
  try {
    // Check if deps directory exists as a simple way to verify if deps are installed
    const depsDir = join(appDir, 'deps')
    if (existsSync(depsDir)) {
      _console.info('Dependencies already installed')
      return true
    }

    if (!userInstallChoice) {
      _console.warn('Dependencies not installed but required for next steps')
      _console.info('Installing dependencies now...')

      // Change to the app directory before running the command
      process.chdir(appDir)
      execSync('mix deps.get', { stdio: 'inherit' })
      _console.success('Dependencies installed successfully')
    } else {
      // This shouldn't happen if user chose to install, but just in case
      _console.error('Dependencies should be installed but are missing')
      return false
    }

    return true
  } catch (error) {
    _console.error('Failed to install dependencies:', error)
    return false
  }
}

/**
 * Generates a secret key using mix phx.gen.secret
 * @param appDir Directory of the Phoenix application
 * @param userInstallChoice Whether the user chose to install dependencies
 * @returns The generated secret key
 */
function generateSecretKey(appDir: string, userInstallChoice: boolean): string {
  try {
    // Ensure dependencies are installed before generating secret key
    if (!ensureDependenciesInstalled(appDir, userInstallChoice)) {
      _console.warn('Could not ensure dependencies are installed')
      return 'REPLACE_WITH_GENERATED_SECRET_KEY'
    }

    // Change to the app directory before running the command
    process.chdir(appDir)
    const secretKey = execSync('mix phx.gen.secret', { encoding: 'utf8' }).trim()
    _console.info('Generated secret key for Phoenix application')
    return secretKey
  } catch (error) {
    _console.warn('Failed to generate secret key:', error)
    return 'REPLACE_WITH_GENERATED_SECRET_KEY'
  }
}

/**
 * Creates .env and .env.example files for Phoenix application
 * @param appDir Directory of the Phoenix application
 * @param userInstallChoice Whether the user chose to install dependencies
 */
function createEnvFiles(appDir: string, userInstallChoice: boolean): void {
  try {
    // Generate a secret key, passing the app directory and user install choice
    const secretKey = generateSecretKey(appDir, userInstallChoice)

    // Ensure secretKey doesn't contain compilation logs
    const cleanSecretKey = secretKey.trim().split('\n').pop() || 'REPLACE_WITH_GENERATED_SECRET_KEY'

    // For .env.example, don't include the actual secret key
    const envExampleContent = 'SECRET_KEY_BASE='

    // For .env, include the actual secret key
    const envContent = `SECRET_KEY_BASE=${cleanSecretKey}`

    const envExamplePath = join(appDir, '.env.example')
    const envPath = join(appDir, '.env')

    // Create .env.example file
    writeFileSync(envExamplePath, envExampleContent)
    _console.success(`Created ${envExamplePath}`)

    // Create .env file
    writeFileSync(envPath, envContent)
    _console.success(`Created ${envPath}`)
  } catch (error) {
    _console.error('Failed to create .env files:', error)
  }
}

/**
 * Generates release files for Phoenix application
 * @param appDir Directory of the Phoenix application
 */
function generateReleaseFiles(appDir: string): void {
  try {
    _console.info('Generating release files with mix phx.gen.release...')
    process.chdir(appDir)
    execSync('mix phx.gen.release', { stdio: 'inherit' })
    _console.success('Release files generated successfully')
  } catch (error) {
    _console.error('Failed to generate release files:', error)
  }
}

/**
 * Copies and customizes moon.yml from template to the app directory
 * @param appDir Directory of the Phoenix application
 * @param appName Name of the application
 * @param appDescription Description of the application
 */
function copyMoonYml(appDir: string, appName: string, appDescription: string): void {
  try {
    _console.info('Copying moon.yml configuration...')

    // Path to template moon.yml
    const templatePath = join(process.cwd(), '..', '..', 'template-phoenix', 'moon.yml')

    // Read template content
    let moonYmlContent = readFileSync(templatePath, 'utf8')

    // Replace variables
    moonYmlContent = moonYmlContent
      .replace(/{{ package_name }}/g, appName)
      .replace(/{{ package_description }}/g, appDescription)

    // Write to app directory
    const destPath = join(appDir, 'moon.yml')
    writeFileSync(destPath, moonYmlContent)

    _console.success(`Created ${destPath}`)
  } catch (error) {
    _console.error('Failed to copy moon.yml:', error)
  }
}

/**
 * Generate a Phoenix Framework application
 * @param options Options for generating a Phoenix Framework application
 * @returns Promise that resolves when the application is generated
 */
export async function generatePhoenixApp(options: PhoenixGenerateOptions): Promise<void> {
  const { appName, appDescription, database, assets, force, install } = options

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

    // Generate release files
    if (install) {
      generateReleaseFiles(appDir)
    } else {
      _console.info('Skipping release generation as dependencies were not installed')
    }

    // Create .env and .env.example files - pass user's install choice
    createEnvFiles(appDir, install)

    // Copy and customize moon.yml
    copyMoonYml(appDir, appName, appDescription)

    _console.success(`Phoenix Framework application ${appName} generated successfully!`)

    if (!install) {
      _console.success(`Run \`moon ${appName}:deps\` to install dependencies.`)
    }
  } catch (error) {
    _console.error('Failed to generate Phoenix Framework application:', error)
    throw error
  }
}
