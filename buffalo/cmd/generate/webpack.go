// Copyright © 2016 Mark Bates <mark@markbates.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package generate

import (
	"fmt"
	"os/exec"

	"github.com/markbates/gentronics"
	"github.com/spf13/cobra"
)

var publicLogo = &gentronics.RemoteFile{
	File:       gentronics.NewFile("public/assets/images/logo.svg", ""),
	RemotePath: "https://raw.githubusercontent.com/gobuffalo/buffalo/master/logo.svg",
}

var assetsLogo = &gentronics.RemoteFile{
	File:       gentronics.NewFile("assets/images/logo.svg", ""),
	RemotePath: "https://raw.githubusercontent.com/gobuffalo/buffalo/master/logo.svg",
}

// WebpackCmd generates a new actions/resource file and a stub test.
var WebpackCmd = &cobra.Command{
	Use:   "webpack",
	Short: "Generates a webpack asset pipeline.",
	RunE: func(cmd *cobra.Command, args []string) error {
		data := gentronics.Data{
			"withWebpack": true,
		}
		return NewWebpackGenerator(data).Run(".", data)
	},
}

// NewWebpackGenerator generates a new actions/resource file and a stub test.
func NewWebpackGenerator(data gentronics.Data) *gentronics.Generator {
	g := gentronics.New()

	should := func(data gentronics.Data) bool {
		if b, ok := data["withWebpack"]; ok {
			return b.(bool)
		}
		return false
	}

	// if we're not using web pack save the logo and return
	if !should(data) {
		g.Add(publicLogo)
		return g
	}

	// if there's no npm, return!
	_, err := exec.LookPath("npm")
	if err != nil {
		fmt.Println("Could not find npm/node. Skipping webpack generation.")
		g.Add(publicLogo)
		return g
	}

	g.Should = should
	g.Add(assetsLogo)
	g.Add(gentronics.NewFile("webpack.config.js", nWebpack))
	g.Add(gentronics.NewFile("public/assets/.gitignore", ""))
	g.Add(gentronics.NewFile("assets/js/application.js", wApplicationJS))
	g.Add(gentronics.NewFile("assets/css/application.scss", wApplicationCSS))
	c := gentronics.NewCommand(exec.Command("npm", "install", "webpack", "-g"))
	g.Add(c)
	c = gentronics.NewCommand(exec.Command("npm", "init", "-y"))
	g.Add(c)

	modules := []string{"webpack", "sass-loader", "css-loader", "style-loader", "node-sass",
		"babel-loader", "extract-text-webpack-plugin", "babel", "babel-core", "url-loader", "file-loader",
		"jquery", "bootstrap", "path", "font-awesome", "npm-install-webpack-plugin", "jquery-ujs",
		"copy-webpack-plugin",
	}
	args := []string{"install", "--save"}
	args = append(args, modules...)
	g.Add(gentronics.NewCommand(exec.Command("npm", args...)))
	return g
}

var nWebpack = `var webpack = require("webpack");
var CopyWebpackPlugin = require('copy-webpack-plugin');
var ExtractTextPlugin = require("extract-text-webpack-plugin");

module.exports = {
  entry: [
    "./assets/js/application.js",
    "./assets/css/application.scss",
    "./node_modules/jquery-ujs/src/rails.js"
  ],
  output: {
    filename: "application.js",
    path: "./public/assets"
  },
  plugins: [
    new webpack.ProvidePlugin({
      $: "jquery",
      jQuery: "jquery"
    }),
    new ExtractTextPlugin("application.css"),
    new CopyWebpackPlugin([{
      from: "./assets",
      to: ""
    }], {
      ignore: [
        "css/*",
        "js/*",
      ]
    })
  ],
  module: {
    loaders: [{
      test: /\.jsx?$/,
      loader: "babel",
      exclude: /node_modules/
    }, {
      test: /\.scss$/,
      loader: ExtractTextPlugin.extract(
        "style",
        "css?sourceMap!sass?sourceMap"
      )
    }, {
      test: /\.woff(\?v=\d+\.\d+\.\d+)?$/,
      loader: "url?limit=10000&mimetype=application/font-woff"
    }, {
      test: /\.woff2(\?v=\d+\.\d+\.\d+)?$/,
      loader: "url?limit=10000&mimetype=application/font-woff"
    }, {
      test: /\.ttf(\?v=\d+\.\d+\.\d+)?$/,
      loader: "url?limit=10000&mimetype=application/octet-stream"
    }, {
      test: /\.eot(\?v=\d+\.\d+\.\d+)?$/,
      loader: "file"
    }, {
      test: /\.svg(\?v=\d+\.\d+\.\d+)?$/,
      loader: "url?limit=10000&mimetype=image/svg+xml"
    }]
  }
};
`

const wApplicationJS = `require("bootstrap/dist/js/bootstrap.js");

$(() => {

});`
const wApplicationCSS = `@import "~bootstrap/dist/css/bootstrap.css";
@import "~font-awesome/css/font-awesome.css";
`
