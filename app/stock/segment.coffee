d3 = require 'd3'
KLine = require './kline'

colors = ["#000", "#000", "#f00", "#080", "#f00", "#080"]

class KLineSegment
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.segment
    @_ui = @root._ui

  init: ->

  update: (data, datasel, dataset) ->
    svg = @_ui.svg
    @_ui.svg.select("g.segment").remove()
    if not datasel.Segment
      return
    if not datasel.Segment.Data
      return
    g = @_ui.svg.append("g")
      .attr("class", "segment")

    x = @_ui.x
    y = @_ui.y
    left = @root._left
    size = @root.options.size

    sdata = datasel.Segment.Data
    dataset = []
    last = {}
    for d in sdata
      if not d.oI
        d.oI = t.I for t in datasel.Typing.Line when t.Time == d.Time
      d.i = d.oI - left

      if d.i >= 0 and d.i <= size
        if last.i < 0 or last.i > size
          dataset.push last
        dataset.push d
      else if last.i >= 0 and last.i <= size
        dataset.push d
      last = d

    dataset = KLine.filter sdata, data
    color = (d, i) -> colors[d.Type] || colors[0]
    g.selectAll('circle')
      .data(dataset)
      .enter()
      .append('circle')
      .attr('cx', (d) -> x d.i)
      .attr('cy', (d) -> y d.Price)
      .attr('r', 6)
      .style("stroke", color)
      .style("fill", (d,i) -> if d.Case1 then color(d,i) else '#fff')
      .on('mouseover', (d,i) -> console.log(d,i))

KLine.register_plugin 'segment', KLineSegment
