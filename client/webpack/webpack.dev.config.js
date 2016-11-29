var path = require('path')
var webpack = require('webpack')

var config = require('./webpack.base.config.js')

config.profile = false
config.devtool = 'inline-source-map'

config.entry.unshift('react-hot-loader/patch')

config.plugins = config.plugins.concat([
  new webpack.NoErrorsPlugin(),
  new webpack.DllReferencePlugin({
    context: path.join(__dirname, '../app'),
    manifest: require('../.dll/vendor-manifest.json')
  })
])

module.exports = config
