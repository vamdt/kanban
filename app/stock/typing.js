import plugin from './plugin';
import { extend, filter } from './util';

class Typing {
  constructor(root) {
    this.root = root;
    this.options = extend({}, this.root.options.typing);
    this._ui = this.root._ui;
  }

  init() {
  }

  update(data, datasel) {
    let tdata = datasel.Typing.Data;
    if (this.root.param('ntyping')) {
      tdata = false;
    }
    const dataset = filter(tdata, data);

    const style = {
      fill: this._ui.tColor,
    };
    this._ui.circle(dataset, 'typing', style);
  }
}

plugin.register('typing', Typing);
