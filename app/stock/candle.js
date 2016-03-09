import Mas from './mas';
import plugin from './plugin';
import { extend } from './util';

const defaults = {
  width: 2,
};

class Candle {
  constructor(root) {
    this.root = root;
    this.options = extend({}, this.root.options.candle, defaults);
    this.root.options.candle = this.options;
    this._ui = this.root._ui;
  }

  init() {
    const svg = this.root._ui.svg;
    this.options.width = +this.options.width || 4;
    this.root.options.size = Math.floor(this.root.options.width / (1 + this.options.width));
    const mas = new Mas(this.root, svg, this.root._ui.y, (d) => d.close);
    mas.init();
    this.root.add_plugin_obj(mas);
  }

  update(data) {
    const svg = this._ui.svg;

    if (this.root.param('nc')) {
      svg.selectAll('.candle').remove();
      return;
    }

    const kColor = (d, i) => this._ui.kColor(d, i, data);
    const x = this._ui.x;
    const y = this._ui.y;
    const candleWidth = this.options.width;
    const dispatch = this.root.dispatch;

    const rect = svg.selectAll('rect.candle')
      .data(this.root.param('ocl') ? [] : data);

    function mover(d, i) { dispatch.tip(this, 'k', d, i); }
    rect
      .enter()
      .append('rect')
      .attr('class', 'candle')
      .attr('width', candleWidth)
      .on('mouseover', mover);

    rect
      .exit()
      .transition()
      .remove();

    rect
      .transition()
      .attr('x', (d, i) => x(i) - candleWidth / 2)
      .attr('y', d => y(Math.max(d.open, d.close)))
      .attr('height', d => Math.max(0.5, Math.abs(y(d.open) - y(d.close))))
      .attr('stroke', kColor)
      .attr('fill', kColor);

    const line = svg.selectAll('line.candle')
      .data(data);

    line
      .enter()
      .append('line')
      .attr('class', 'candle')
      .style('stroke-width', '1')
      .on('mouseover', mover);

    line
      .exit()
      .transition()
      .remove();

    line
      .transition()
      .style('stroke', kColor)
      .attr('x1', (d, i) => x(i) - (d.Low === d.High ? candleWidth / 2 : 0))
      .attr('y1', d => y(d.High))
      .attr('x2', (d, i) => x(i) + (d.Low === d.High ? candleWidth / 2 : 0))
      .attr('y2', d => y(d.Low));

    const opacity = this.root.param('opacity');
    if (opacity) {
      svg.selectAll('.candle')
        .transition()
        .style('opacity', opacity);
    }
  }

}

plugin.register('candle', Candle);
