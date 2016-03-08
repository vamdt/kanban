const plugins = {};

export default {
  register(name, clazz) {
    plugins[name] = clazz;
  },
  every(fn) {
    for (const name in plugins) {
      if (plugins.hasOwnProperty(name)) {
        fn(name, plugins[name]);
      }
    }
  },
};
