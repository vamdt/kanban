import d3 from 'd3';
import KLine from './kline';

const bColor = (d) => (d.MACD > 0) ? '#f00' : '#080';

const formatValue = d3.format(',.3f');

function fmtMacd(d) {
  return formatValue(d / 1000);
}

function dragmove() {
  d3.select(this)
    .attr('transform', `translate(0, ${d3.event.y})`);
}

class KLineMacd {
  constructor(root) {
    this.root = root;
    this.options = KLine.extend({}, this.root.options.macd);
  }

  init() {
    const height = this.height = 50;
    const rsvg = this.root._ui.svg;
    const w = +d3.select(rsvg.node().parentNode).attr('width');
    this.w = w;
    if (this.options.container) {
      const margin = this.root.options.margin;
      const container = d3.select(this.options.container);
      this.svg = container.append('svg')
        .attr('width', w)
        .attr('height', height)
        .append('g')
        .attr('transform', `translate(${margin.left},0)`);
    } else {
      const drag = d3.behavior.drag().on('drag', dragmove);
      this.svg = rsvg.append('g')
        .call(drag);
    }
  }

  hide() {
    this.svg
      .transition()
      .style('display', 'none');
  }

  show() {
    this.svg
      .transition()
      .style('display', '');
  }

  updateAxis(data) {
    if (this.root.param('nmacd')) {
      this.hide();
      return;
    }

    const svg = this.svg;
    let axis = svg.select('#macd_y_axis');
    if (axis.empty()) {
      axis = svg.append('g')
        .attr('class', 'y axis')
        .attr('id', 'macd_y_axis');
      this.y = d3.scale.linear()
        .range([this.height, 0]);
      this.yAxis = d3.svg.axis()
        .scale(this.y)
        .orient('left')
        .ticks(4)
        .tickSize(-this.w)
        .tickFormat(fmtMacd);
    }
    const min = d3.min(data, d => Math.min(d.DIFF, d.DEA, d.MACD));
    const max = d3.max(data, d => Math.max(d.DIFF, d.DEA, d.MACD));
    this.y.domain([min, max]);
    axis.call(this.yAxis);
    this.show();
  }

  update(data) {
    if (this.root.param('nmacd')) {
      this.hide();
      return;
    }
    this.show();

    const x = this.root._ui.x;
    const y = this.y;
    const candleWidth = this.root.options.candle.width;
    const svg = this.svg;

    const rect = svg.selectAll('rect.macd')
      .data(data);

    rect
      .exit()
      .transition()
      .remove();
    rect
      .enter()
      .append('rect')
      .attr('class', 'macd')
      .attr('width', candleWidth);
    rect
      .transition()
      .attr('x', (d, i) => x(i) - candleWidth / 2)
      .attr('y', (d) => y(Math.max(d.MACD, 0)))
      .attr('height', (d) => Math.abs(y(0) - y(d.MACD)))
      .attr('stroke', bColor)
      .attr('fill', bColor);

    let ldiff = svg.select('path#diff');
    if (ldiff.empty()) {
      ldiff = svg.append('path')
        .attr('id', 'diff')
        .style('fill', 'none')
        .style('stroke', 'silver')
        .style('stroke-width', '1');
    }
    ldiff.data([data]);

    let line = d3.svg.line()
      .x((d, i) => x(i))
      .y((d) => y(d.DIFF));

    ldiff.attr('d', line);

    let ldea = svg.select('path#dea');
    if (ldea.empty()) {
      ldea = svg.append('path')
        .attr('id', 'dea')
        .style('fill', 'none')
        .style('stroke', 'gold')
        .style('stroke-width', '1');
    }
    ldea.data([data]);

    line = d3.svg.line()
      .x((d, i) => x(i))
      .y((d) => y(d.DEA));

    ldea.attr('d', line);
  }
}

KLine.register_plugin('macd', KLineMacd);
