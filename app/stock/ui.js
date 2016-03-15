import d3 from 'd3';
import util from './util';

const color = {
  up: '#f00',
  down: '#080',
  eq: '#000',
};
color._up = color.up;
color._down = color.down;
color._eq = color.eq;

const formatValue = d3.format(',.2f');
function fmtCent(d) { return formatValue(d / 100); }

function xAxisTickFormat(root) {
  let _prevTick = 0;
  const _Ymd = d3.time.format('%Y-%m-%d');
  const _md = d3.time.format('%m-%d');
  const _dHM = d3.time.format('%d %H:%M');
  const _HM = d3.time.format('%H:%M');
  const _M = d3.time.format(':%M');
  return function format(i) {
    const data = root.data();
    if (typeof data[i] === 'undefined') {
      return 'F';
    }
    const date = data[i].date;
    const prevDate = data.length > _prevTick ? data[_prevTick].date : date;
    _prevTick = i;

    if (i === 0) {
      return _Ymd(date);
    }

    if (date === prevDate) {
      return '';
    }

    if (date.getYear() !== prevDate.getYear()) {
      return _Ymd(date);
    }
    if (date.getMonth() !== prevDate.getMonth()) {
      return _md(date);
    }
    if (date.getDay() !== prevDate.getDay()) {
      return _dHM(date);
    }
    if (date.getHours() !== prevDate.getHours()) {
      return _HM(date);
    }
    return _M(date);
  };
}

function zoomFn(root) {
  return function zoomfn() {
    const n = d3.event.scale;
    this.zs = this.zs || n;
    const o = this.zs;
    this.zs = n;

    const x1 = d3.event.translate[0];
    this.zx = this.zx || x1;
    const x0 = this.zx;

    let nsize = root.options.size;
    let nleft = root._left;
    if (n < o) {
      nsize = parseInt(nsize * 1.1, 10);
    } else if (n > o) {
      nsize = parseInt(nsize * 0.9, 10);
    } else {
      if (Math.abs(Math.abs(x1) - Math.abs(x0)) < 2) {
        return;
      }
      this.zx = x1;

      if (x0 > x1) {
        nleft = nleft + Math.max(20, parseInt(nsize * 0.05, 10));
      } else if (x0 < x1) {
        nleft = nleft - Math.max(20, parseInt(nsize * 0.05, 10));
      } else {
        return;
      }
    }
    root.update_size(nsize, nleft);
    root.delay_draw();
  };
}

export default class KUI {
  constructor(root) {
    this.root = root;
    this.dispatch = root.dispatch;
    root.dispatch.on('param.ui', () => this.updateColor());
    root.dispatch.on('resize.ui', () => this.resize());
  }

  init() {
    const root = this.root;
    const options = this.root.options;
    const margin = options.margin;

    this.container = d3.select(options.container || 'body');
    const container = this.container;
    container.html('');
    let width = parseInt(container.style('width'), 10);
    if (width < 1) {
      width = util.w();
    }
    let height = parseInt(container.style('height'), 10);
    if (height < 1) {
      height = util.h() - 80;
    }

    options.width = width - margin.left - margin.right;
    options.height = height - margin.top - margin.bottom;
    width = options.width;
    height = options.height;

    const svg = this.svg = container.append('svg')
      .attr('width', width + margin.left + margin.right)
      .attr('height', height + margin.top + margin.bottom)
      .append('g')
      .attr('transform', `translate(${margin.left},${margin.top})`);

    this.x = d3.scale.linear()
      .range([0, width]);

    this.y = d3.scale.linear()
      .range([height, 0]);

    this.xAxis = d3.svg.axis()
      .scale(this.x)
      .orient('bottom')
      .tickSize(-height, 0)
      .tickFormat(xAxisTickFormat(root));

    this.yAxis = d3.svg.axis()
      .scale(this.y)
      .orient('left')
      .ticks(6)
      .tickSize(-width)
      .tickFormat(fmtCent);

    svg.append('g')
      .attr('class', 'x axis')
      .attr('transform', `translate(0, ${height})`)
      .call(this.xAxis);

    svg.append('g')
      .attr('class', 'y axis')
      .call(this.yAxis);

    const zoom = d3.behavior.zoom()
      .on('zoom', zoomFn(root));

    svg
      .call(zoom);
    svg.append('rect')
      .attr('class', 'pane')
      .attr('width', width)
      .attr('height', height);
  }

  update(data) {
    this.x.domain([0, data.length - 1]);
    this.y.domain([d3.min(data, (d) => d.Low) * 0.99, d3.max(data, (d) => d.High)]);

    this.svg.select('.x.axis')
      .call(this.xAxis);
    this.svg.select('.y.axis')
      .call(this.yAxis);
  }

  resize() {
    const options = this.root.options;
    const width = options.width;
    const height = options.height;
    const margin = options.margin;
    this.container.select('svg')
        .attr('width', width + margin.left + margin.right)
        .attr('height', height + margin.top + margin.bottom);
    this.x.range([0, width]);
    this.y.range([height, 0]);
    this.xAxis.tickSize(-height, 0);
    this.yAxis.tickSize(-width);
    this.svg.select('.x.axis')
      .attr('transform', `translate(0, ${height})`);
    this.svg.select('rect.pane')
      .attr('width', width)
      .attr('height', height);
    this.root.delay_draw();
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
