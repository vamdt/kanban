import KLine from './kline';

class KLineTyping {
  constructor(root) {
    this.root = root;
    this.options = KLine.extend({}, this.root.options.typing);
    this._ui = this.root._ui;
  }

  init() {
  }

  update(data, datasel) {
    let tdata = datasel.Typing.Data;
    if (this.root.param('ntyping')) {
      tdata = false;
    }
    const dataset = KLine.filter(tdata, data);

    const style = {
      fill: this.root.tColor,
    };
    this._ui.circle(dataset, 'typing', style);
  }
}

KLine.register_plugin('typing', KLineTyping);
