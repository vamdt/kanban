import KLine from './kline';

class KLineSegment {
  constructor(root) {
    this.root = root;
    this.options = KLine.extend({}, this.root.options.segment);
    this._ui = this.root._ui;
  }

  init() {
  }

  update(data, datasel) {
    const sdata = datasel.Segment.Data;
    const dataset = KLine.filter(sdata, data);

    const style = {
      stroke: this.root.tColor,
      fill: (d) => d.Case1 ? this.root.tColor(d) : '#fff',
    };
    this._ui.circle(dataset, 'segment', style);
  }
}

KLine.register_plugin('segment', KLineSegment);
