d3 = require 'd3'
KLine = require './kline'

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
      .on('mouseover', (d, i) -> dispatch.tip @, 'typing', d, i)

    circle.exit().transition().remove()

    [eq, up, down] = [KLine.color.eq, KLine.color.up, KLine.color.down]
    colors = [eq, eq, up, down, up, down]

    rsize = @root.param('typing_circle_size') || 3

    circle
      .transition()
      .attr('r', rsize)
      .attr('cx', (d) -> x d.i)
      .attr('cy', (d) -> y d.Price)
      .style("fill", (d,i) -> colors[d.Type] || colors[0])

KLine.register_plugin 'typing', KLineTyping
