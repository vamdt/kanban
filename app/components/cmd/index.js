import sugg from './sugg';
import unwatch from './unwatch';
import watch from './watch';
import show from './show';

const kanpan = [
  'begin',
  'hub',
];

const handlers = {
  sugg,
  unwatch,
  watch,
  show,
};

export default {
  ready() {
    const bind = (e, func) => {
      this.$on(e, (opt) => {
        func.apply(this, opt);
      });
    };
    for (const e in handlers) {
      if (handlers.hasOwnProperty(e)) {
        bind(e, handlers[e]);
      }
    }

    kanpan.forEach((e) => {
      this.$on(e, (opt) => {
        opt.unshift(e);
        this.$root.$broadcast('kline_cmd', opt);
      });
    });
  },
};
