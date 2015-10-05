d3 = require 'd3'
KLine = require './kline'
defaults =
  '5d':
    interval: 5
    color: '#abc'
  '10d':
    interval: 10
    color: 'red'

class KLineMas
  constructor: (@root, @svg, @y, @d) ->
    @options = KLine.extend {}, @root.options.mas, defaults
    @root.options.mas = @options
    @y = @y || @root._ui.y
    @d = @d || (d) -> d.close

  init: ->
    svg = @svg
    for n,ma of @options
      interval = +ma.interval

      line = svg.select("path.line.mas.ma#{interval}")
      if line.empty()
        line = svg.append("path")
          .attr("class", "line mas ma#{interval}")
      line.attr("clip-path", "url(#clip)")
        .style("stroke", ma.color)
        .style("stroke-width", "1")

  update: (data) ->
    svg = @svg
    for n,ma of @options
      interval = +ma.interval
      e = svg.select("path.line.mas.ma#{interval}")
      e.data([data])
      @drawMA(data, interval, e)

  drawMA: (data, interval, element) ->
    x = @root._ui.x
    y = @y
    left = @root._left
    data = @root._data
    dfn = @d
    mean = (d, i) ->
      l = left + i - interval - 1
      l = Math.max(l, 0)
      d3.mean data[l..(left+i)], dfn

    line = d3.svg.line()
      .x((d, i) -> x i)
      .y((d, i) -> y mean d, i)

    element.attr("d", line)

module.exports = KLineMas
