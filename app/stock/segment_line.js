import plugin from './plugin';
import { extend, filter } from './util';

class KLineSegmentLine {
  constructor(root) {
    this.root = root;
    this.options = extend({}, this.root.options.segment);
    this._ui = this.root._ui;
  }

  init() {
  }

  update(data, datasel) {
    let dname = 'Line';
    const handcraft = this.root.param('handcraft');
    if (handcraft && this.root.param('k') !== '1') {
      dname = 'HCLine';
    }
    const dataset = filter(datasel.Segment[dname], data);
    const up = 4;
    dataset.forEach((d, j) => {
      if (d.hasOwnProperty('MACD')) {
        return;
      }
      if (d.i < 0 || d.ei < 0) {
        return;
      }

      let mup = 0;
      let mdown = 0;
      for (let i = d.i; i < d.ei && i < data.length; i++) {
        if (data[i].MACD > 0) {
          mup += data[i].MACD;
        } else {
          mdown += data[i].MACD;
        }
      }
      dataset[j].MACD = d.Type === up ? mup : mdown;
    });

    this._ui.line(dataset, 'segment_line');
    if (handcraft) {
      const begin = datasel.begin || 0;
      this._ui.lineno(dataset, begin, 'segment_line');
    } else {
      this._ui.svg.selectAll('text.segment_line').remove();
    }
  }
}

plugin.register('segment_line', KLineSegmentLine);
