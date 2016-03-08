import d3 from 'd3';
import KLineMas from './mas';
import plugin from './plugin';
import { extend } from './util';

const formatValue = d3.format(',d');
const fmtVolume = (d) => formatValue(d / 100);

function dragmove() {
  d3.select(this).attr('transform', `translate(0, ${d3.event.y})`);
}

class KLineVolume {
  constructor(root) {
    this.root = root;
    this.options = extend({}, this.root.options.volume);
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
      const drag = d3.behavior.drag()
        .on('drag', dragmove);
      this.svg = rsvg.append('g')
        .call(drag);
    }
    this.y = d3.scale.linear()
      .range([height, 0]);
    const mas = new KLineMas(this.root, this.svg, this.y, (d) => d.volume);
    mas.init();
    this.root.add_plugin_obj(mas);
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
    if (this.root.param('nvolume')) {
      this.hide();
      return;
    }

    const svg = this.svg;
    let axis = svg.select('#volume_y_axis');
    if (axis.empty()) {
      axis = svg.append('g')
        .attr('class', 'y axis')
        .attr('id', 'volume_y_axis');
      this.yAxis = d3.svg.axis()
        .scale(this.y)
        .orient('left')
        .ticks(4)
        .tickSize(-this.w)
        .tickFormat(fmtVolume);
    }
    this.y.domain([0, d3.max(data, (d) => d.volume)]);
    axis.call(this.yAxis);
    this.show();
  }

  update(data) {
    if (this.root.param('nvolume')) {
      this.hide();
      return;
    }

    const kColor = (d, i) => this.root._ui.kColor(d, i, data);
    const x = this.root._ui.x;
    const y = this.y;
    const height = this.height;
    const candleWidth = this.root.options.candle.width;
    const svg = this.svg;

    const rect = svg.selectAll('rect.volume')
      .data(data);

    rect
      .exit()
      .transition()
      .remove();

    rect
      .enter()
      .append('rect')
      .attr('class', 'volume')
      .attr('width', candleWidth);

    rect
      .transition()
      .attr('x', (d, i) => x(i) - candleWidth / 2)
      .attr('y', (d) => y(d.volume))
      .attr('height', (d) => height - y(d.volume))
      .attr('stroke', kColor)
      .attr('fill', kColor);
    this.show();
  }
}

plugin.register('volume', KLineVolume);
