import Vue from 'vue';
import KLine from '../stock/webpack';
import config from './config';

Vue.directive('kanpan', {
  deep: true,

  bind() {
    this.vm.$on('kline_cmd', (opt) => {
      if (!this.kl) {
        return;
      }
      this.kl.cmd.apply(this.kl, opt);
    });

    window.addEventListener('resize', () => {
      if (this.kl) {
        this.kl.resize();
      }
    });

    window.addEventListener('keyup', (e) => {
      if (e.target.tagName === 'INPUT') {
        return;
      }
      if (!this.kl) {
        return;
      }
      const handles = {
        49: 'nmacd',
        50: 'nvolume',
        51: 'nc',
        52: 'nmas',
      };
      const name = handles[e.keyCode] || false;
      if (!name) {
        return;
      }
      const param = {};
      param[name] = !this.kl.param(name);
      config.update(param);
      this.kl.param(param);
    });

    window.addEventListener('keydown', (e) => {
      if (e.target.tagName === 'INPUT') {
        return;
      }
      if (!this.kl) {
        return;
      }
      const kl = this.kl;
      const handles = {
        35() {
          kl.move('end');
        },
        36() {
          kl.move('home');
        },
        37() {
          kl.move('left');
        },
        39() {
          kl.move('right');
        },
      };
      if (handles.hasOwnProperty(e.keyCode)) {
        handles[e.keyCode]();
      }
    });
  },

  update(value) {
    if (!value || !value.s) {
      return;
    }
    const settings = config.load() || {};
    const params = JSON.parse(JSON.stringify(value));
    for (const k in params) {
      if (params.hasOwnProperty(k)) {
        settings[k] = params[k];
      }
    }
    if (this.kl) {
      this.kl.param(settings);
      return;
    }

    this.kl = new KLine({
      container: this.el,
    });
    this.kl.dispatch.on('nameChange', (v) => {
      this.vm.sname = v;
    });

    this.vm.$nextTick(() => {
      this.kl.param(settings);
    });
  },
});
