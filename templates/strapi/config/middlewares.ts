export default ({ env }) => {
  const awsBucket = env('AWS_BUCKET')
  const awsRegion = env('AWS_REGION', 'ap-southeast-1')
  const s3AssetsUrl = awsBucket ? `${awsBucket}.s3.${awsRegion}.amazonaws.com` : undefined

  return [
    'strapi::logger',
    'strapi::errors',
    {
      name: 'strapi::security',
      config: {
        contentSecurityPolicy: {
          useDefaults: true,
          directives: {
            'connect-src': ["'self'", 'https:'],
            'script-src': ["'self'", "'unsafe-inline'"],
            'img-src': ["'self'", 'data:', 'blob:', s3AssetsUrl],
            'media-src': ["'self'", 'data:', 'blob:', s3AssetsUrl],
            upgradeInsecureRequests: null,
          },
        },
      },
    },
    {
      name: 'strapi::cors',
      config: {
        headers: '*',
        origin: [env('STRAPI_ALLOWED_ORIGIN', '*')],
      },
    },
    'strapi::query',
    'strapi::body',
    'strapi::session',
    'strapi::favicon',
    'strapi::public',
  ]
}
