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
    data = datasel.Typing.Data
    @_ui.svg.select("g.typing").remove()
    g = @_ui.svg.append("g")
      .attr("class", "typing")

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

    g.selectAll('circle')
      .data(dataset)
      .enter()
      .append('circle')
      .attr('cx', (d) -> x d.i)
      .attr('cy', (d) -> y d.Price)
      .attr('r', 2)
      .style("fill", (d,i) -> colors[d.Type] || colors[0])
      .on('mouseover', (d,i) -> console.log(d,i))

KLine.register_plugin 'typing', KLineTyping
