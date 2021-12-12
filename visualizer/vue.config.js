module.exports = {
  devServer: {
    proxy: {
      '/ws': {
        target: 'http://localhost:4000',
        ws: true,
        changeOrigin: true
      }
    }
  }
}
