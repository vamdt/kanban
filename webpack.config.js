var path = require("path");
var webpack = require("webpack");
var ExtractTextPlugin = require('extract-text-webpack-plugin');
module.exports = {
	cache: true,
  entry: {
    main: "./app/scripts/main.coffee"
  },
	output: {
		path: path.join(__dirname, ".tmp"),
		filename: "[name].js",
		chunkFilename: "[chunkhash].js"
	},
  devtool: 'source-map',
  resolve: {
    extensions: ['', '.coffee', '.js'],
    modulesDirectories: ['node_modules', 'scripts', 'styles', 'images', 'fonts']
  },
	module: {
    noParse: /\.min\.js/,
		loaders: [
			{
        test: /\.css$/,
        //loader: ExtractTextPlugin.extract('style-loader', 'css-loader?modules&importLoaders=1&localIdentName=[name]__[local]___[hash:base64:5]!postcss-loader')
        loader: 'style-loader!css-loader?modules&importLoaders=1&localIdentName=[name]__[local]___[hash:base64:5]!postcss-loader'
      },

			// required for bootstrap icons
			{ test: /\.woff$/,   loader: "url-loader?prefix=font/&limit=5000&mimetype=application/font-woff" },
			{ test: /\.ttf$/,    loader: "file-loader?prefix=font/" },
			{ test: /\.eot$/,    loader: "file-loader?prefix=font/" },
			{ test: /\.svg$/,    loader: "file-loader?prefix=font/" },

      { test: /\.png$/, loader: "url-loader?mimetype=image/png" },
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
	plugins: [
    //new ExtractTextPlugin("main.css"),
		new webpack.ProvidePlugin({
			// Automtically detect jQuery and $ as free var in modules
			// and inject the jquery library
			// This is required by many jquery plugins
			jQuery: "jquery",
			$: "jquery"
		})
	]
};
