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
    const x = this._ui.x;
    const y = this._ui.y;
    const dispatch = this.root.dispatch;
    const sdata = datasel.Segment.Data;
    const dataset = KLine.filter(sdata, data);
    const [eq, up, down] = [KLine.color.eq, KLine.color.up, KLine.color.down];
    const colors = [eq, eq, up, down, up, down];
    const color = (d) => colors[d.Type] || colors[0];
    const c = this._ui.svg.selectAll('circle.segment')
      .data(dataset);

    function mover(d, i) { dispatch.tip(this, 'segment', d, i); }
    c
      .enter()
      .append('circle')
      .attr('class', 'segment')
      .on('mouseover', mover);

    c
      .exit()
      .transition()
      .remove();

    const rsize = this.root.param('segment_circle_size') || 3;
    c
      .transition()
      .attr('r', rsize)
      .attr('cx', (d) => x(d.i))
      .attr('cy', (d) => y(d.Price))
      .style('stroke', color)
      .style('fill', (d, i) => d.Case1 ? color(d, i) : '#fff');
  }
}

KLine.register_plugin('segment', KLineSegment);
