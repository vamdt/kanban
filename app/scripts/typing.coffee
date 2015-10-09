d3 = require 'd3'
KLine = require './kline'

colors = ["#000", "#000", "#f00", "#080", "#f00", "#080"]

class KLineTyping
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.typing
    @_ui = @root._ui

  init: ->

  update: (data, datasel, dataset) ->
    svg = @_ui.svg
    svg.selectAll('circle').remove()
    data = datasel.Typing
    g = @_ui.svg.select("g.typing")
    if !data or data.length < 1
      if !g.empty()
        g.remove()
      return
    if g.empty()
      g = @_ui.svg.append("g")
        .attr("class", "typing")

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

    svg.selectAll('circle')
      .data(data)
      .enter()
      .append('g')
      .append('circle')
      .attr('cx', line.x())
      .attr('cy', line.y())
      .attr('r', 3)
      .style("fill", (d,i) -> colors[d.Type] || colors[0])
      .on('mouseover', (d,i) -> console.log(d,i))

KLine.register_plugin 'typing', KLineTyping
