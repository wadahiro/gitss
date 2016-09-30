var webpack = require('webpack')

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
      // {
      //   test: /\.woff(\?v=\d+\.\d+\.\d+)?$/,
      //   loader: "url-loader?mimetype=application/font-woff"
      // },
      // {
      //   test: /\.woff2(\?v=\d+\.\d+\.\d+)?$/,
      //   loader: "url-loader?mimetype=application/font-woff"
      // },
      // {
      //   test: /\.ttf(\?v=\d+\.\d+\.\d+)?$/,
      //   loader: "url-loader?mimetype=application/font-woff"
      // },
      // {
      //   test: /\.eot(\?v=\d+\.\d+\.\d+)?$/,
      //   loader: "url-loader?mimetype=application/font-woff"
      // },
      // {
      //   test: /\.svg(\?v=\d+\.\d+\.\d+)?$/,
      //   loader: "url-loader?mimetype=image/svg+xml"
      // },
      {
        test: /\.css$/,
        loader: 'style-loader!css-loader?modules',
      },
      {
        test: /\.js(x?)$/,
        exclude: [/node_modules/],
        loaders: ['babel-loader?cacheDirectory=true']
      },
      {
        test: /\.ts(x?)$/,
        exclude: [/node_modules/],
        loaders: ['babel-loader?cacheDirectory=true', 'awesome-typescript-loader?forkChecker=true']
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
    new webpack.DllReferencePlugin({
      context: path.join(__dirname, '../app'),
      manifest: require('../.dll/vendor-manifest.json')
    })
  ],
  cache: true
}
