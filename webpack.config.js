var path = require("path");
var webpack = require("webpack");
var ExtractTextPlugin = require('extract-text-webpack-plugin');
module.exports = {
  cache: true,
  context: __dirname + "/app",
  entry: {
    main: "./main.coffee"
  },
  output: {
    path: path.join(__dirname, "dist", "static"),
    filename: "[name].js",
    chunkFilename: "[chunkhash].js"
  },
  devtool: 'source-map',
  resolve: {
    extensions: ['', '.coffee', '.js'],
    modulesDirectories: ['node_modules', 'scripts', 'styles', 'images', 'fonts']
  },
  module: {
    loaders: [
    {
      test: /\.css$/,
      //loader: ExtractTextPlugin.extract('style-loader', 'css-loader?modules&importLoaders=1&localIdentName=[name]__[local]___[hash:base64:5]!postcss-loader')
      loader: 'style!css!postcss'
    },
    { test: /\.html$/, loader: "file?name=[name].[ext]" },

    // required for bootstrap icons
    { test: /\.woff$/,   loader: "url-loader?prefix=font/&limit=5000&mimetype=application/font-woff" },
    { test: /\.ttf$/,    loader: "file-loader?prefix=font/" },
    { test: /\.eot$/,    loader: "file-loader?prefix=font/" },
    { test: /\.svg$/,    loader: "file-loader?prefix=font/" },

    { test: /\.png$/, loader: "url-loader?limit=100000&mimetype=image/png" },
    { test: /\.gif$/, loader: "file-loader" },
    { test: /\.jpg$/, loader: "file-loader" },
    { test: /\.vue$/, loader: "vue" },
    { test: /\.coffee$/, loader: "coffee-loader?sourceMap" }
    ]
  },
  postcss: [
    require('autoprefixer'),
    require('postcss-color-rebeccapurple')
  ],
  //plugins: [
    ////new ExtractTextPlugin("[name].css"),
    //new webpack.ProvidePlugin({
      //// Automtically detect jQuery and $ as free var in modules
      //// and inject the jquery library
      //// This is required by many jquery plugins
      //jQuery: "jquery",
      //$: "jquery"
    //})
  //],
  devServer: {
    port: 9001,
    progress: true,
    hot: true,
    quiet: false,
    noInfo: false,
    stats: { colors: true },
    historyApiFallback: false,
    contentBase: 'app',
    proxy: {
      '*': 'http://127.0.0.1:3002'
    },
  }
};
