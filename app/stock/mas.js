import d3 from 'd3';

const defaults = [
  {
    interval: 5,
    color: 'silver',
  },
  {
    interval: 10,
    color: 'gray',
  },
];

let ref = 0;

function defaultDfn(d) {
  return d.close;
}

export default class KLineMas {
  constructor(root, svg, y, dfn) {
    this.root = root;
    this.svg = svg;
    this.y = y || root._ui.y;
    this.d = dfn || defaultDfn;
    this.ref = ++ref;
  }

  init() {
  }

  update(data) {
    if (this.root.param('nmas')) {
      this.svg.selectAll('path.mas').remove();
      return;
    }
    const mas = this.root.param('mas') || defaults;
    const color = d3.scale.category20();
    const dispatch = this.root.dispatch;
    function mover(d, i) { dispatch.tip(this, 'mas', d, i); }
    mas.forEach((ma) => {
      const interval = +ma.interval;
      const id = `ma${interval}-${this.ref}`;
      let e = this.svg.select(`path#${id}`);
      if (e.empty()) {
        e = this.svg.append('path')
          .attr('class', 'mas')
          .attr('id', id)
          .attr('stroke', ma.color || color(interval))
          .attr('stroke-width', '1')
          .attr('fill', 'none');
      }

      e
        .data([data])
        .on('mouseover', mover);
      this.drawMA(interval, e);
    });
  }

  drawMA(interval, element) {
    const x = this.root._ui.x;
    const y = this.y;
    const left = this.root._left;
    const data = this.root._data;
    const dfn = this.d;
    function mean(d, i) {
      let l = left + i - interval - 1;
      l = Math.max(l, 0);
      return d3.mean(data.slice(l, left + i + 1), dfn);
    }

    const line = d3.svg.line()
      .x((d, i) => x(i))
      .y((d, i) => y(mean(d, i)));

    element
      .transition()
      .attr('d', line);
  }

}
