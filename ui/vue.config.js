module.exports = {
  lintOnSave: false,
  productionSourceMap: false,
  outputDir: "dist/files",
  devServer: {
    //disableHostCheck: true,
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:48080',
        changeOrigin: true,
      },
      '/extend': {
        target: 'http://127.0.0.1:48080',
        changeOrigin: true,
      },
    },
  },
};
