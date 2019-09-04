let path = require('path')
let glob = require('glob')
//配置pages多页面获取当前文件夹下的html和js
function getEntry(globPath) {
	let entries = {},
		basename, tmp, pathname;

	glob.sync(globPath).forEach(function (entry) {
		basename = path.basename(entry, path.extname(entry));
		// console.log(entry)
		tmp = entry.split('/').splice(-3);

		pathname = basename; // 正确输出js和html的路径

		// console.log(pathname)
		entries[pathname] = {
			entry: 'src/' + tmp[0] + '/' + tmp[1] + '/' + tmp[1] + '.js',
			template: 'src/' + tmp[0] + '/' + tmp[1] + '/' + tmp[2],
			title: tmp[2],
			filename: tmp[2]
		};
	});
	return entries;
}

let pages = getEntry('./src/pages/**?/*.html');



module.exports = {
	lintOnSave: false, //禁用eslint
	publicPath: '/static',
	outputDir: "../views",
	runtimeCompiler: true,
	productionSourceMap: false,
	pages: {
		index: {
			// page 的入口
			entry: 'src/pages/index/main.js',
			filename: 'index.html',
		},
		login: {
			// page 的入口
			entry: 'src/pages/login/main.js',

			filename: 'login.html',
			minify: false
		}
	}
}