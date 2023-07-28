'use strict';

var svgpath = require('svgpath');

exports.name = 'iconvg';
exports.type = 'perItem';
exports.active = true;
exports.description = 'try to make it compatible with iconvg';
exports.params = {};
exports.fn = () => {
  return {
    element: {
      enter: (node, parentNode) => {
        // Replaces all arcs with Bézier curves.
        if (node.name === 'path' && node.attributes.d != null) {
          var transformed = svgpath(node.attributes.d).unarc().toString()
          node.attributes.d = transformed;
        }
        // Remove fill-rule="evenodd" from all paths.
        if (node.name === 'g' && node.attributes['fill-rule'] === 'evenodd') {
          delete node.attributes['fill-rule'];
        }
      },
    },
  };
}
