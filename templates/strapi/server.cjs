#!/usr/bin/env node --no-warnings

/* @ref: https://docs.strapi.io/dev-docs/deployment */

const strapi = require('@strapi/strapi')
const path = require('node:path');

strapi.createStrapi({
  distDir: path.resolve(__dirname, './dist'),
  serveAdminPanel: true,
  autoReload: true,
}).start()
