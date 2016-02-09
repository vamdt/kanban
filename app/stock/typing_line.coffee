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

KLine.register_plugin 'typing_line', KLineTypingLine
