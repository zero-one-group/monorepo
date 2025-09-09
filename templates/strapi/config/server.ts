export default ({ env }) => ({
  host: env('HOST', '0.0.0.0'),
  port: env.int('PORT', {{ port_number }}),
  app: {
    keys: env.array('APP_KEYS'),
  },
})
