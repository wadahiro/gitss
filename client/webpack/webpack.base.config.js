var webpack = require('webpack')
var ExtractTextPlugin = require('extract-text-webpack-plugin')

var path = require('path')
var objectAssign = require('object-assign')

var NODE_ENV = process.env.NODE_ENV

var env = {
  production: NODE_ENV === 'production',
  staging: NODE_ENV === 'staging',
  test: NODE_ENV === 'test',
  development: NODE_ENV === 'development' || typeof NODE_ENV === 'undefined'
}

objectAssign(env, {
  build: (env.production || env.staging)
})

module.exports = {
  target: 'web',
  entry: ['babel-polyfill', path.join(__dirname, '../index.tsx')],
  output: {
    path: path.join(__dirname, '../../assets'),
    filename: 'js/bundle.js'
  },
  module: {
    loaders: [
      {
        test: /\.woff(2)?(\?v=[0-9]\.[0-9]\.[0-9])?$/,
        loader: 'url-loader?limit=10000&mimetype=application/font-woff&publicPath=../&name=./css/[hash].[ext]'
      },
      {
        test: /\.(ttf|eot|svg)(\?v=[0-9]\.[0-9]\.[0-9])?$/,
        loader: 'file-loader?publicPath=../&name=./css/[hash].[ext]'
      },
      {
        test: /\.css$/,
        loader: ExtractTextPlugin.extract({
          fallbackLoader: 'style-loader',
          loader: 'css-loader'
        })
      },
      {
        test: /\.js(x?)$/,
        exclude: [/node_modules/],
        loaders: ['babel-loader?cacheDirectory=true']
      },
      {
        test: /\.ts(x?)$/,
        exclude: [/node_modules/],
        loaders: ['babel-loader?cacheDirectory=true', 'ts-loader?transpileOnly=false']
      },
      //   {
      //     test: /\.css$/,
      //     loader: "style!css"
      //   }
    ]
  },
  resolve: {
    extensions: ['.tsx', '.ts', '.js', '.jsx']
  },
  plugins: [
    new ExtractTextPlugin('css/style.css')
  ],
  cache: true
}
