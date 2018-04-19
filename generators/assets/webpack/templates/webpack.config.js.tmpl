const Webpack = require("webpack");
const Glob = require("glob");
const path = require("path");
const CopyWebpackPlugin = require("copy-webpack-plugin");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const ManifestPlugin = require("webpack-manifest-plugin");
const CleanObsoleteChunks = require('webpack-clean-obsolete-chunks');
const UglifyJsPlugin = require("uglifyjs-webpack-plugin");
const LiveReloadPlugin = require('webpack-livereload-plugin');

const configurator = {
  entries: function(){
    var entries = {
      application: [
        './node_modules/jquery-ujs/src/rails.js',
        './assets/css/application.scss',
      ],
    }
    
    Glob.sync("./assets/*/*.*").forEach((entry) => {
      if (entry === './assets/css/application.scss') {
        return
      }
      
      let key = entry.replace(/(\.\/assets\/(js|css|go)\/)|\.(js|s[ac]ss|go)/g, '')
      if(key.startsWith("_") || (/(js|s[ac]ss|go)$/i).test(entry) == false) {
        return
      }
      
      if( entries[key] == null) {
        entries[key] = [entry]
        return
      } 
      
      entries[key].push(entry)
    })

    return entries
  },

  plugins() {
    var plugins = [
      new CleanObsoleteChunks(),
      new Webpack.ProvidePlugin({$: "jquery",jQuery: "jquery"}),
      new MiniCssExtractPlugin({filename: "[name].[contenthash].css"}),
      new CopyWebpackPlugin([{from: "./assets",to: ""}], {copyUnmodified: true,ignore: ["css/**", "js/**"] }),
      new Webpack.LoaderOptionsPlugin({minimize: true,debug: false}),
      new ManifestPlugin({fileName: "manifest.json"})
    ];

    return plugins
  },

  moduleOptions: function() {
    return {
      rules: [
        {
          test: /\.s[ac]ss$/,
          use: [
            MiniCssExtractPlugin.loader,
            { loader: "css-loader", options: {sourceMap: true}},
            { loader: "sass-loader", options: {sourceMap: true}}
          ]
        },
        { test: /\.jsx?$/,loader: "babel-loader",exclude: /node_modules/ },
        { test: /\.(woff|woff2|ttf|svg)(\?v=\d+\.\d+\.\d+)?$/,use: "url-loader"},
        { test: /\.eot(\?v=\d+\.\d+\.\d+)?$/,use: "file-loader" },
        { test: require.resolve("jquery"),use: "expose-loader?jQuery!expose-loader?$"},
        { test: /\.go$/, use: "gopherjs-loader"}
      ]
    }
  },

  buildConfig: function(){
    const env = process.env.NODE_ENV || "development";
    
    var config = {
      mode: env,
      entry: configurator.entries(),
      output: {filename: "[name].[hash].js", path: `${__dirname}/public/assets`},
      plugins: configurator.plugins(),
      module: configurator.moduleOptions()
    }

    if( env === "development" ){
      config.plugins.push(new LiveReloadPlugin({appendScriptTag: true}))
      return config
    }

    const uglifier = new UglifyJsPlugin({
      uglifyOptions: {
        beautify: false,
        mangle: {keep_fnames: true},
        output: {comments: false},
        compress: {}
      }
    })

    config.optimization = {
      minimizer: [uglifier]
    }

    return config
  }
}

module.exports = configurator.buildConfig()