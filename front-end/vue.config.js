const webpack = require('webpack');
module.exports = {
    configureWebpack: {
        plugins: [
            new webpack.DefinePlugin({
                'process.env': {
                    'SERVER_PATH': JSON.stringify(process.env.SERVER_PATH),
                    'HTTP_PATH': JSON.stringify(process.env.HTTP_PATH),
                }
            })
        ]
    }
};
