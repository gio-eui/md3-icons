//
const iconvg = require('./svgo.iconvg.js');

// svgo.config.js
module.exports = {
  js2svg: {
    indent: 2,
    pretty: true,
  },
  plugins: [
    // convert all path data into absolute path data
    {
      fn: iconvg.fn,
      name: iconvg.name,
      type: iconvg.type,
      active: iconvg.active,
      description: iconvg.description,
    },
    // Default plugins
    {
      name: 'preset-default',
      params: {
        overrides: {
          mergePaths: false,
          convertPathData:false,
          removeViewBox: true,
        },
      },
    },
  ],
};
