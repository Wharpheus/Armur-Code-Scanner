const pluginJs = require('@eslint/js');

module.exports = {
  plugins: {
    js: pluginJs.configs.recommended
  },
  rules: {
    "no-unused-vars": "error",
    "no-unreachable": "error",
    "no-constant-condition": "error",
    "no-unused-expressions": "error",
    "no-unused-private-class-members": "error",
    "no-useless-assignment": "error"
  }
};
