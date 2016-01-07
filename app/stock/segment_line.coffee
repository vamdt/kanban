d3 = require 'd3'
KLine = require './kline'

colors = ["#000", "#000", "#f00", "#080", "#f00", "#080"]

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

    path.attr("d", line)

    # TODO opt for last line of path
    if dataset.length and dataset[dataset.length-1] == ldata[ldata.length-1]
      d = dataset[dataset.length-1]
      if not d.ei
        return
      l = g.append("line")
        .style("fill", "none")
        .style("stroke", '#abc')
        .style("stroke-width", "2")
        .attr("x1", x(d.i)).attr("x2", x(d.ei))
      if d.Type is 5
        l.attr("y1", y(d.High)).attr("y2", y(d.Low))
      else
        l.attr("y1", y(d.Low)).attr("y2", y(d.High))

KLine.register_plugin 'segment_line', KLineSegmentLine
