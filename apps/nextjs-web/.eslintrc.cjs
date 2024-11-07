/** @type {import('@types/eslint').Linter.BaseConfig} */
module.exports = {
  extends: ['next/core-web-vitals', 'eslint:recommended', 'plugin:@typescript-eslint/recommended'],
  plugins: ['@typescript-eslint', 'simple-import-sort'],
  ignorePatterns: ['node_modules', '.next', '.cache', 'build', 'dist', 'out'],
  settings: {
    next: {
      rootDir: __dirname,
    },
  },
  rules: {
    'react/jsx-key': 'off',
    'import/no-anonymous-default-export': 'off',
    'no-console': [
      'warn',
      {
        allow: ['error', 'info'],
      },
    ],
    '@next/next/no-img-element': 'off',
    '@next/next/no-html-link-for-pages': 'off',
    '@typescript-eslint/no-non-null-assertion': 'off',
    '@typescript-eslint/no-explicit-any': 'off',
    '@typescript-eslint/no-var-requires': 'off',
    '@typescript-eslint/no-unused-vars': [
      'warn',
      {
        argsIgnorePattern: '^_',
        varsIgnorePattern: '^_',
        caughtErrorsIgnorePattern: '^_',
      },
    ],
  },
  overrides: [
    {
      files: ['*.js', '*.jsx', '*.ts', '*.tsx'],
      rules: {
        'simple-import-sort/exports': 'warn',
        'simple-import-sort/imports': [
          'error',
          {
            groups: [
              // Dotenv and config packages first.
              ['^dotenv', '^(#/env)(/.*|$)'],
              // React related packages come after.
              ['^react', '^next(?!/?$)', '^@?\\w'],
              // Internal packages.
              ['^(#|components)(/.*|$)'],
              // Side effect imports.
              ['^\\u0000'],
              // Parent imports. Put `..` last.
              ['^\\.\\.(?!/?$)', '^\\.\\./?$'],
              // Other relative imports. Put same-folder imports and `.` last.
              ['^\\./(?=.*/)(?!/?$)', '^\\.(?!/?$)', '^\\./?$'],
              // Style imports.
              ['^.+\\.?(css)$'],
            ],
          },
        ],
      },
    },
  ],
  globals: {
    React: true,
    JSX: true,
  },
}
