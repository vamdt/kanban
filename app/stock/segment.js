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
      .style('stroke', this.root.tColor)
      .style('fill', (d) => d.Case1 ? this.root.tColor(d) : '#fff');
  }
}

KLine.register_plugin('segment', KLineSegment);
