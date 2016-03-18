import d3 from 'd3';
import { extend, mergeData } from './util';
import Plugin from './plugin';
import KUI from './ui';
import Data from './data';

const defaults = {
  container: 'body',
  margin: {
    top: 20,
    right: 50,
    bottom: 30,
    left: 50,
  },
};

export default class KLine {
  constructor(options) {
    this.options = options;
    const ev = [
      'resize',
      'param',
      'tip',
      'cmd',
      'redraw',
      'nameChange',
      'uiInit',
    ];
    this.dispatch = d3.dispatch(...ev);
    this.options = extend({}, this.options, defaults);
    this._data = [];
    this._ui = new KUI(this);
    this.io = new Data();
    this._left = 0;
    this._max_left = 0;
    this.plugins = [];
    this._param = {};
    this._cache = {};
    this.options.size = +this.options.size || 100;
    this.bindEvent();
  }

  datalen() {
    return this._data ? this._data.length : 0;
  }

  updateSize(_size, _left) {
    let size = _size || this.options.size || 10;
    size = Math.max(size, 10);
    this.options.size = size;

    const left = _left || this._left;

    const atrightedge = this._left === this._max_left;

    this._max_left = Math.max(0, this.datalen() - 1 - this.options.size);

    this.options.size = this.datalen() - 1 - this._max_left;

    if (atrightedge && left === this._left) {
      this._left = this._max_left;
    } else {
      this._left = Math.min(this._max_left, Math.max(0, left));
    }
  }

  cmd(...args) {
    this.dispatch.cmd.apply(this, args);
  }

  data() {
    this._data = this._data || [];
    return this._data.slice(this._left, this._left + this.options.size + 1);
  }

  setData(_data) {
    if (!_data || !_data.id) {
      return;
    }
    const sid = _data.id;
    this._cache[sid] = mergeData(this._cache[sid], _data);
  }

  selData() {
    const s = this.param('s');
    this._dataset = this._cache[s] || false;
    const data = this._dataset || {};

    const levels = {
      1: 'm1s',
      5: 'm5s',
      30: 'm30s',
      day: 'days',
      week: 'weeks',
      month: 'months',
    };

    const k = this.param('k');
    this._datasel = data[levels[k]] || data.days || {};
    this._data = this._datasel.data;
    this.updateSize();
    if (data.Name !== this.title) {
      this.title = data.Name;
      this.dispatch.nameChange(data.Name);
    }
  }

  param(p) {
    const t = typeof p;

    if (t === 'string') {
      return this._param[p];
    }

    if (t === 'object') {
      const o = {};
      Object.keys(p).forEach((k) => {
        const v = p[k];
        if (this._param[k] !== v) {
          o[k] = this._param[k];
        }
        this._param[k] = v;
      });
      return this.dispatch.param(o);
    }

    return this._param;
  }

  bindEvent() {
    const redraw = (data) => {
      if (!data || !data.id) {
        return;
      }
      this.setData(data);
      this.selData();
      this.delayDraw();
    };

    this.dispatch.on('param.core', (o) => {
      if (o.hasOwnProperty('s')) {
        this.io.subscribe(this.param('s'), redraw);
      }
      this.dispatch.redraw();
    });

    this.dispatch.on('redraw.core', () => {
      this.selData();
      this.delayDraw();
    });

    this.dispatch.on('uiInit.core', () => {
      this.initPlugins();
    });
  }

  addPluginObj(plugin) {
    this.plugins.push(plugin);
  }

  initPlugins() {
    this.plugins = [];
    Plugin.eachDo((name, C) => {
      const plugin = new C(this);
      if (plugin.init() !== false) {
        this.addPluginObj(plugin);
      }
    });
  }

  resize(...args) {
    this.dispatch.resize(...args);
  }

  move(dir) {
    let left = 0;
    switch (dir) {
      case 'left':
        left = this._left - 100;
        break;
      case 'right':
        left = this._left + 100;
        break;
      case 'home':
        left = -1;
        break;
      case 'end':
        left = 100000;
        break;
      default:
        return;
    }
    this.updateSize(0, +left);
    this.delayDraw();
  }

  delayDraw() {
    if (this.__need_draw) {
      this.__need_draw++;
      return;
    }
    this.__need_draw = 1;
    d3.timer(() => {
      this.draw();
      return true;
    });
  }

  draw() {
    const data = this.data();
    this._ui.update(data);
    this.update(data, this._datasel, this._dataset);
    this.__need_draw = 0;
  }

  update(data, datasel, dataset) {
    this.plugins.forEach((plugin) => {
      if (plugin.updateAxis) {
        plugin.updateAxis(data, datasel, dataset);
      }
      plugin.update(data, datasel, dataset);
    });
  }

  stop() {
    this.io.close();
  }

  notification(id, msg) {
    if (!window.Notification) {
      return;
    }

    const show = () => {
      const config = {
        body: msg,
        dir: 'auto',
      };
      this.notif = new Notification(id, config);
    };

    if (Notification.permission === 'granted') {
      show();
    } else if (Notification.permission !== 'denied') {
      Notification.requestPermission((permission) => {
        if (permission === 'granted') {
          show();
        }
      });
    }
  }
}
