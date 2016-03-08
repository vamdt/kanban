import d3 from 'd3';
import plugin from './plugin';
import { extend, filter } from './util';

class KLineHub {
  constructor(root) {
    this.root = root;
    this.options = extend({}, this.root.options.hub);
    this._ui = this.root._ui;
  }

  init() {
  }

  update(data, datasel, dataset) {
    const ksel = this.root.param('k');
    const levels = [
      {
        level: '1',
        name: 'm1s',
      },
      {
        level: '5',
        name: 'm5s',
      },
      {
        level: '30',
        name: 'm30s',
      },
      {
        level: 'day',
        name: 'days',
      },
      {
        level: 'week',
        name: 'weeks',
      },
      {
        level: 'month',
        name: 'months',
      },
    ];
    const handcraft = this.root.param('handcraft');
    const dname = handcraft ? 'Data' : 'HCData';
    let skip = false;
    levels.forEach(({ level, name }) => {
      const hubdata = (!skip && dataset[name]) ? dataset[name].Hub[dname] : false;
      this.draw(level, hubdata, data);
      if (level === ksel) {
        skip = true;
      }
    });
  }

  draw(k, data, kdata) {
    const cls = `hub-${k}`;
    const x = this._ui.x;
    const y = this._ui.y;
    const dispatch = this.root.dispatch;
    const dataset = filter(data, kdata);
    const rect = this._ui.svg.selectAll(`rect.${cls}`)
      .data(dataset);
    const hover = {
      stroke: 'green',
      'stroke-width': '1',
      'stroke-opacity': '1',
    };
    const hout = {
      stroke: 'steelblue',
      'stroke-width': '1',
      'stroke-opacity': '.3',
    };

    rect
      .enter()
      .append('rect')
      .attr('class', cls)
      .attr('fill', 'transparent')
      .style(hout)
      .on('mouseover.stroke', () => d3.select(this).style(hover))
      .on('mouseover.tip', (d, i) => dispatch.tip(this, 'hub', d, i))
      .on('mouseout.stroke', () => d3.select(this).style(hout));

    rect
      .exit()
      .transition()
      .remove();

    rect
      .transition()
      .attr('x', (d) => x(d.i))
      .attr('y', (d) => y(d.High))
      .attr('width', (d) => Math.max(0.5, x(d.ei) - x(d.i)))
      .attr('height', (d) => Math.max(0.5, y(d.Low) - y(d.High)));

    const text = this._ui.svg.selectAll(`text.${cls}`)
      .data(dataset);
    text
      .enter()
      .append('text')
      .attr('class', cls)
      .attr('fill', 'black');

    text
      .exit()
      .transition()
      .remove();

    text
      .transition()
      .attr('x', (d) => x(d.i))
      .attr('y', (d) => y(d.High) + 10)
      .text(k);
  }
}

plugin.register('hub', KLineHub);
