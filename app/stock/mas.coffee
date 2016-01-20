d3 = require 'd3'
KLine = require './kline'
defaults = [ {interval: 5, color: 'silver'}, {interval: 10, color: 'gray'} ]

class KLineMas
  constructor: (@root, @svg, @y, @d) ->
    @y = @y || @root._ui.y
    @d = @d || (d) -> d.close

  init: ->

  update: (data) ->
    svg = @svg
    if @root.param('nmas')
      svg.selectAll("path.mas").remove()
      return
    mas = @root.param('mas') || defaults
    color = d3.scale.category20()
    dispatch = @root.dispatch
    for ma in mas
      interval = +ma.interval
      e = svg.select("path#ma#{interval}")
      if e.empty()
        e = svg.append("path")
          .attr("class", "line mas")
          .attr("id", "ma#{interval}")
          .style("stroke", ma.color||color interval)
          .style("stroke-width", "1")
      e.data([data])
        .on('mouseover', (d, i) -> dispatch.tip @, 'mas', d, i)
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

    element
      .transition()
      .attr("d", line)

module.exports = KLineMas
