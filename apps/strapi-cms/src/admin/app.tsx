/* @see: https://docs.strapi.io/dev-docs/admin-panel-customization */

import type { StrapiApp } from '@strapi/strapi/admin'
import { consola } from 'consola'
import brandIcon from '../assets/brand-icon.svg'

export default {
  config: {
    head: {
      favicon: brandIcon,
    },
    auth: {
      logo: brandIcon,
    },
    menu: {
      logo: brandIcon,
    },
    locales: ['id'],
    tutorials: false, // Disable video tutorials
    notifications: {
      releases: false, // Disable notifications about new Strapi releases
    },
  },
  bootstrap(app: StrapiApp) {
    consola.debug(app)
  },
}
