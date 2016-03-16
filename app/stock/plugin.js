const plugins = {};

export default {
  register(name, clazz) {
    plugins[name] = clazz;
  },
  eachDo(fn) {
    for (const name in plugins) {
      if (plugins.hasOwnProperty(name)) {
        fn(name, plugins[name]);
      }
    }
  },
};
