d3 = require 'd3'
KLine = require './kline'

colors = ["#000", "#000", "#f00", "#080", "#f00", "#080"]

class KLineTyping
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.typing
    @_ui = @root._ui

  init: ->

  update: (data, datasel, dataset) ->
    @_ui.svg.select("g#typing").remove()
    return unless datasel.Typing.Data
    g = @_ui.svg.append("g")
      .attr("id", "typing")

    x = @_ui.x
    y = @_ui.y

    tdata = datasel.Typing.Data
    dataset = KLine.filter tdata, data
    dispatch = @root.dispatch

    g.selectAll('circle')
      .data(dataset)
      .enter()
      .append('circle')
      .attr('cx', (d) -> x d.i)
      .attr('cy', (d) -> y d.Price)
      .attr('r', 2)
      .style("fill", (d,i) -> colors[d.Type] || colors[0])
      .on('mouseover', (d, i) -> dispatch.tip @, 'typing', d, i)

KLine.register_plugin 'typing', KLineTyping
