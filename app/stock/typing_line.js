import plugin from './plugin';
import { extend, filter } from './util';

class TypingLine {
  constructor(root) {
    this.root = root;
    this.options = extend({}, this.root.options.typing);
    this._ui = this.root._ui;
  }

  init() {
  }

  update(data, datasel) {
    const dset = filter(datasel.Typing.Line, data);
    const style = {
      'stroke-dasharray': '7 7',
      stroke: '#abc',
    };
    this._ui.line(dset, 'typing_line', style);
  }
}

plugin.register('typing_line', TypingLine);
