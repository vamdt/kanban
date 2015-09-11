var path = require("path");
var webpack = require("webpack");
module.exports = {
	cache: true,
	entry: "./app/scripts/main.coffee",
	output: {
		path: path.join(__dirname, "dist"),
		publicPath: "dist/static/scripts",
		filename: "[name].js",
		chunkFilename: "[chunkhash].js"
	},
  resolve: {
    extensions: ['', '.coffee', '.js']
  },
	module: {
		loaders: [
			// required to write "require('./style.css')"
			{ test: /\.css$/,    loader: "style-loader!css-loader!autoprefixer-loader" },

			// required for bootstrap icons
			{ test: /\.woff$/,   loader: "url-loader?prefix=font/&limit=5000&mimetype=application/font-woff" },
			{ test: /\.ttf$/,    loader: "file-loader?prefix=font/" },
			{ test: /\.eot$/,    loader: "file-loader?prefix=font/" },
			{ test: /\.svg$/,    loader: "file-loader?prefix=font/" },

      { test: /\.coffee$/, loader: "coffee-loader" }
		]
	},
	plugins: [
		new webpack.ProvidePlugin({
			// Automtically detect jQuery and $ as free var in modules
			// and inject the jquery library
			// This is required by many jquery plugins
			jQuery: "jquery",
			$: "jquery"
		})
	]
};
