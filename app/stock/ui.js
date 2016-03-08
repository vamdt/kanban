import d3 from 'd3';

const color = {
  up: '#f00',
  down: '#080',
  eq: '#000',
};
color._up = color.up;
color._down = color.down;
color._eq = color.eq;

export default class KUI {
  constructor(root) {
    this.root = root;
    this.dispatch = root.dispatch;
    root.dispatch.on('param.ui', () => this.updateColor());
  }

  updateColor() {
    const rcolor = this.root.param('color');
    if (!rcolor) {
      return;
    }
    ['up', 'down', 'eq'].forEach((n) => {
      if (rcolor.hasOwnProperty(n)) {
        color[n] = rcolor[n] || color[`_${n}`];
      }
    });
  }

  path(dataset, id, style1 = {}) {
    const style = style1;
    style.fill = style.fill || 'none';
    const x = this.x;
    const y = this.y;
    let path = this.svg.select(`path#${id}`);
    if (path.empty()) {
      path = this.svg.append('path')
        .attr('id', id);
    }
    path
      .style(style)
      .data([dataset]);

    const line = d3.svg.line()
      .x((d) => x(d.i))
      .y((d) => y(d.Price));

    path
      .transition()
      .attr('d', line);
  }

  line(dataset, clazz, style1 = {}) {
    const dispatch = this.dispatch;

    const style = style1;
    style.fill = style.fill || 'none';
    style.strokeWidth = style.strokeWidth || '1';
    const x = this.x;
    const y = this.y;

    const line = this.svg.selectAll(`line.${clazz}`)
      .data(dataset);

    line
      .exit()
      .transition()
      .remove();

    line
      .enter()
      .append('line')
      .attr('class', clazz)
      .on('mouseover.tip', function mover(d, i) { dispatch.tip(this, clazz, d, i); })
      .style(style);

    const up = 4;
    const down = 5;

    line
      .transition()
      .attr('x1', (d) => x(d.i))
      .attr('y1', (d) => y(d.Type === up ? d.Low : d.High))
      .attr('x2', (d) => x(d.ei))
      .attr('y2', (d) => y(d.Type === down ? d.Low : d.High))
      .style('stroke', style.stroke || this.tColor);
  }

  lineno(dataset, begin, clazz, style = {}) {
    const x = this.x;
    const y = this.y;
    const up = 4;
    const down = 5;

    const data = dataset.map((e) => e);
    if (data.length > 0) {
      let d = data[0];
      d = {
        ei: d.i,
        Type: d.Type,
        Low: d.Low,
        High: d.High,
        no: d.no - 1,
      };
      d.Type = (d.Type === up) ? down : up;
      data.unshift(d);
    }

    const text = this.svg.selectAll(`text.${clazz}`)
      .data(data);

    text
      .exit()
      .transition()
      .remove();

    text
      .enter()
      .append('text')
      .attr('class', clazz)
      .style(style);

    const numf = (d) => {
      const n = d.no + 1 - begin;
      return (n > -1) ? n : '';
    };

    text
      .transition()
      .attr('x', (d) => x(d.ei))
      .attr('y', (d) => y(d.Type === down ? d.Low : d.High))
      .text(numf)
      .style('stroke', style.stroke || this.tColor);
  }

  circle(dataset, clazz, style = {}) {
    const x = this.x;
    const y = this.y;
    const dispatch = this.root.dispatch;
    const circle = this.svg.selectAll(`circle.${clazz}`)
      .data(dataset);

    function mover(d, i) { dispatch.tip(this, clazz, d, i); }

    circle
      .exit()
      .transition()
      .remove();

    circle
      .enter()
      .append('circle')
      .attr('class', clazz)
      .on('mouseover.tip', mover);

    const rsize = this.root.param(`${clazz}_circle_size`) || 3;

    circle
      .transition()
      .attr('r', rsize)
      .attr('cx', (d) => x(d.i))
      .attr('cy', (d) => y(d.Price))
      .style(style);
  }

  color(compare) {
    return (...args) => compare(...args) ? color.up : color.down;
  }

  tColor(d) {
    return (d.Type === 2 || d.Type === 4) ? color.up : color.down;
  }

  kColor(d, i, data) {
    if (d.open === d.close) {
      if (i > 0 && data) {
        if (data[i] && data[i - 1]) {
          if (data[i].open >= data[i - 1].close) {
            return color.up;
          }
          if (data[i].open < data[i - 1].close) {
            return color.down;
          }
        }
      }
      return color.eq;
    }

    return (d.open > d.close) ? color.down : color.up;
  }
}
