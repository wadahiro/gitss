var webpack = require('webpack')
var WebpackDevServer = require('webpack-dev-server')
var config = require('./webpack.config')

config.entry.unshift('webpack/hot/only-dev-server')
config.entry.unshift('webpack-dev-server/client?http://localhost:9000')

config.module.loaders = config.module.loaders.map(function (loader) {
  if (loader.loaders) {
    // loader.loaders.unshift('react-hot')
  }
  return loader
})

config.plugins.unshift(new webpack.HotModuleReplacementPlugin())

new WebpackDevServer(webpack(config), {
  contentBase: __dirname + '/../../assets',
  hot: true,
  inline: true,
  historyApiFallback: true,
  stats: { colors: true },
  proxy: {
    '/api/*': 'http://localhost:3000'
  }
}).listen(9000, 'localhost', function (err, result) {
  if (err) {
    console.log(err)
  }

  console.log('Listening at localhost:9000')
})
