import sugg from './sugg';
import unwatch from './unwatch';
import watch from './watch';
import show from './show';
import star from './star';
import unstar from './unstar';
import lucky from './lucky';

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
  lucky,
};

export default {
  ready() {
    Object.keys(handlers).forEach((e) => {
      this.$on(e, handlers[e]);
    });

    kanpan.forEach((e) => {
      this.$on(e, (...opt) => {
        this.$root.$broadcast('kline_cmd', e, ...opt);
      });
    });
  },
};
