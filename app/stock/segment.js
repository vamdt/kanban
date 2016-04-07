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

    const fillFn = (d) => {
      if (d.Case1) {
        return this._ui.tColor(d);
      }
      return '#fff';
    };

    const style = {
      stroke: this._ui.tColor,
      fill: fillFn,
    };
    this._ui.circle(dataset, 'segment', style);
  }
}

plugin.register('segment', Segment);
