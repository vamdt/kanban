d3 = require 'd3'
KLine = require './kline'

class KLineTypingLine
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.typing
    @_ui = @root._ui

  init: ->

  update: (data, datasel, dataset) ->
    svg = @_ui.svg
    @_ui.svg.select("g#typing_line").remove()
    ldata = datasel.Typing.Line
    if not ldata
      return
    g = @_ui.svg.append("g")
      .attr("id", "typing_line")
    path = g
      .append("path")
      .style("fill", "none")
      .style("stroke", '#abc')
      .style("stroke-width", "1")
      .style("stroke-dasharray", "7 7")

    x = @_ui.x
    y = @_ui.y

    dataset = KLine.filter ldata, data
    path.data([dataset])

    line = d3.svg.line()
      .x((d) -> x d.i)
      .y((d) -> y d.Price)

    path.transition().attr("d", line)

    # TODO opt for last line of path
    if dataset.length and dataset[dataset.length-1] == ldata[ldata.length-1]
      d = dataset[dataset.length-1]
      if not d.ei
        return
      l = g.append("line")
        .style("fill", "none")
        .style("stroke", '#abc')
        .style("stroke-width", "1")
        .style("stroke-dasharray", "7 7")
        .attr("x1", x(d.i)).attr("x2", x(d.ei))
      if d.Type is 5
        l.attr("y1", y(d.High)).attr("y2", y(d.Low))
      else
        l.attr("y1", y(d.Low)).attr("y2", y(d.High))

KLine.register_plugin 'typing_line', KLineTypingLine
