const encore = require('@symfony/webpack-encore')

encore
  .setOutputPath('./dist')
  .setPublicPath('/assets')
  .addStyleEntry('css/app', './scss/app.scss')
  .disableSingleRuntimeChunk()
  .cleanupOutputBeforeBuild()
  .copyFiles({ from: './images', to: 'images/[path][name].[hash:8].[ext]' })
  .enableSourceMaps(!encore.isProduction())
  .enableVersioning(encore.isProduction())
  .setManifestKeyPrefix("")
  .enableSassLoader()

module.exports = encore.getWebpackConfig()
