d3 = require 'd3'
KLine = require './kline'

class KLineSegmentLine
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.segment
    @_ui = @root._ui

  init: ->

  update: (data, datasel, dataset) ->
    @_ui.svg.select("g.segment_line").remove()
    if not datasel.Segment
      return
    if not datasel.Segment.Line
      return
    g = @_ui.svg.append("g")
      .attr("class", "segment_line")
    path = g
      .append("path")
      .style("fill", "none")
      .style("stroke", '#abc')
      .style("stroke-width", "2")

    x = @_ui.x
    y = @_ui.y

    ldata = datasel.Segment.Line
    dataset = KLine.filter ldata, data
    path.data([dataset])

    line = d3.svg.line()
      .x((d) -> x d.i)
      .y((d) -> y d.Price)

    path.transition().attr("d", line)

KLine.register_plugin 'segment_line', KLineSegmentLine
