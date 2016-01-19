d3 = require 'd3'
KLine = require './kline'

colors = ["#000", "#000", "#f00", "#080", "#f00", "#080"]

class KLineTyping
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.typing
    @_ui = @root._ui

  init: ->

  update: (data, datasel, dataset) ->
    x = @_ui.x
    y = @_ui.y

    tdata = datasel.Typing.Data
    dataset = KLine.filter tdata, data
    dispatch = @root.dispatch

    circle = @_ui.svg.selectAll('circle.typing')
      .data(dataset)

    circle
      .enter()
      .append('circle')
      .attr('class', 'typing')
      .attr('r', 2)
      .on('mouseover', (d, i) -> dispatch.tip @, 'typing', d, i)

    circle.exit().transition().remove()

    circle
      .transition()
      .attr('cx', (d) -> x d.i)
      .attr('cy', (d) -> y d.Price)
      .style("fill", (d,i) -> colors[d.Type] || colors[0])

KLine.register_plugin 'typing', KLineTyping
