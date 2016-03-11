import sugg from './sugg';
import unwatch from './unwatch';
import watch from './watch';
import show from './show';
import star from './star';
import unstar from './unstar';

const kanpan = [
  'begin',
  'hub',
];

const handlers = {
  sugg,
  unwatch,
  watch,
  show,
  star,
  unstar,
};

export default {
  ready() {
    const bind = (e, func) => {
      this.$on(e, (...opt) => {
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
