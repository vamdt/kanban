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
    const x = this._ui.x;
    const y = this._ui.y;
    let tdata = datasel.Typing.Data;
    if (this.root.param('ntyping')) {
      tdata = false;
    }
    const dataset = KLine.filter(tdata, data);
    const dispatch = this.root.dispatch;
    const circle = this._ui.svg.selectAll('circle.typing')
      .data(dataset);
    function mover(d, i) {
      dispatch.tip(this, 'typing', d, i);
    }
    circle
      .enter()
      .append('circle')
      .attr('class', 'typing')
      .on('mouseover', mover);

    circle
      .exit()
      .transition()
      .remove();

    const [eq, up, down] = [KLine.color.eq, KLine.color.up, KLine.color.down];
    const colors = [eq, eq, up, down, up, down];

    const rsize = this.root.param('typing_circle_size') || 3;

    circle
      .transition()
      .attr('r', rsize)
      .attr('cx', (d) => x(d.i))
      .attr('cy', (d) => y(d.Price))
      .style('fill', (d) => colors[d.Type] || colors[0]);
  }
}

KLine.register_plugin('typing', KLineTyping);
