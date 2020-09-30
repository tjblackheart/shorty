const encore = require('@symfony/webpack-encore')

encore
  .setOutputPath('./dist')
  .setPublicPath('/assets')
  .addStyleEntry('css/app', './scss/app.scss')
  .disableSingleRuntimeChunk()
  .cleanupOutputBeforeBuild()
  .enableSourceMaps(!encore.isProduction())
  .enableVersioning(encore.isProduction())
  .setManifestKeyPrefix("")
  .enableSassLoader()

module.exports = encore.getWebpackConfig()
