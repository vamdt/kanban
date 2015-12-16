d3 = require 'd3'
KLine = require './kline'

colors = ["#000", "#000", "#f00", "#080", "#f00", "#080"]

class KLineSegmentLine
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.segment
    @_ui = @root._ui

  init: ->

  update: (data, datasel, dataset) ->
    svg = @_ui.svg
    @_ui.svg.select("g.segment_line").remove()
    if not datasel.Segment
      return
    if not datasel.Segment.Line
      return
    path = @_ui.svg.append("g")
      .attr("class", "segment_line")
      .append("path")
      .style("fill", "none")
      .style("stroke", '#abc')
      .style("stroke-width", "2")

    x = @_ui.x
    y = @_ui.y
    left = @root._left
    size = @root.options.size

    ldata = datasel.Segment.Line
    dataset = []
    last = {}
    for d in ldata
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
    dataset = KLine.filter ldata, data
    path.data([dataset])

    line = d3.svg.line()
      .x((d) -> x d.i)
      .y((d) -> y d.Price)

    path.attr("d", line)

KLine.register_plugin 'segment_line', KLineSegmentLine
