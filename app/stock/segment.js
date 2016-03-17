import plugin from './plugin';
import { extend, filter } from './util';

class Segment {
  constructor(root) {
    this.root = root;
    this.options = extend({}, this.root.options.segment);
    this._ui = this.root._ui;
  }

  init() {
  }

  update(data, datasel) {
    const segment = datasel.Segment || {};
    const sdata = segment.Data;
    const dataset = filter(sdata, data);

    const style = {
      stroke: this._ui.tColor,
      fill: (d) => d.Case1 ? this._ui.tColor(d) : '#fff',
    };
    this._ui.circle(dataset, 'segment', style);
  }
}

plugin.register('segment', Segment);
