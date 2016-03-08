import KLine from './kline';

class KLineTypingLine {
  constructor(root) {
    this.root = root;
    this.options = KLine.extend({}, this.root.options.typing);
    this._ui = this.root._ui;
  }

  init() {
  }

  update(data, datasel) {
    const dset = KLine.filter(datasel.Typing.Line, data);
    const style = {
      'stroke-dasharray': '7 7',
      stroke: '#abc',
    };
    this._ui.line(dset, 'typing_line', style);
  }
}

KLine.register_plugin('typing_line', KLineTypingLine);
