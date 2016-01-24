d3 = require 'd3'
KLine = require './kline'

colors = ["#000", "#000", "#f00", "#080", "#f00", "#080"]

class KLineSegment
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.segment
    @_ui = @root._ui

  init: ->

  update: (data, datasel, dataset) ->
    x = @_ui.x
    y = @_ui.y

    dispatch = @root.dispatch
    sdata = datasel.Segment.Data
    dataset = KLine.filter sdata, data
    color = (d, i) -> colors[d.Type] || colors[0]

    c = @_ui.svg.selectAll("circle.segment")
      .data(dataset)

    c
      .enter()
      .append('circle')
      .attr("class", "segment")
      .attr('r', 4)
      .on('mouseover', (d, i) -> dispatch.tip @, 'segment', d, i)

    c.exit().transition().remove()

    c.transition()
      .attr('cx', (d) -> x d.i)
      .attr('cy', (d) -> y d.Price)
      .style("stroke", color)
      .style("fill", (d,i) -> if d.Case1 then color(d,i) else '#fff')

KLine.register_plugin 'segment', KLineSegment
