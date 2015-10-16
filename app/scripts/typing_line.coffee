d3 = require 'd3'
KLine = require './kline'

colors = ["#000", "#000", "#f00", "#080", "#f00", "#080"]

class KLineTypingLine
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.typing
    @_ui = @root._ui

  init: ->

  update: (data, datasel, dataset) ->
    svg = @_ui.svg
    data = datasel.Line
    g = @_ui.svg.select("g.typing_line")
    if !data or data.length < 1
      if !g.empty()
        g.remove()
      return
    if g.empty()
      g = @_ui.svg.append("g")
        .attr("class", "typing_line")

    path = g.select("path")
    if path.empty()
      path = g.append("path")
    path.style("fill", "none")
      .style("stroke", "blue")
      .style("stroke-width", "1")

    x = @_ui.x
    y = @_ui.y
    left = @root._left
    size = @root.options.size

    dataset = []
    last = {}
    for d in data
      d.i = d.I - left

      if d.i >= 0 and d.i <= size
        if last.i < 0 or last.i > size
          dataset.push last
        dataset.push d
      else if last.i >= 0 and last.i <= size
        dataset.push d
      last = d
    path.data([dataset])

    line = d3.svg.line()
      .x((d) -> x d.i)
      .y((d) -> y d.Price)

    path.attr("d", line)

KLine.register_plugin 'typing_line', KLineTypingLine