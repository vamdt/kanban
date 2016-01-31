d3 = require 'd3'
KLine = require './kline'

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

    [eq, up, down] = [KLine.color.eq, KLine.color.up, KLine.color.down]
    colors = [eq, eq, up, down, up, down]
    color = (d, i) -> colors[d.Type] || colors[0]

    c = @_ui.svg.selectAll("circle.segment")
      .data(dataset)

    c
      .enter()
      .append('circle')
      .attr("class", "segment")
      .on('mouseover', (d, i) -> dispatch.tip @, 'segment', d, i)

    c.exit().transition().remove()

    rsize = @root.param('segment_circle_size') || 3

    c.transition()
      .attr('r', rsize)
      .attr('cx', (d) -> x d.i)
      .attr('cy', (d) -> y d.Price)
      .style("stroke", color)
      .style("fill", (d,i) -> if d.Case1 then color(d,i) else '#fff')

KLine.register_plugin 'segment', KLineSegment
