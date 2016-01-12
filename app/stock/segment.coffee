d3 = require 'd3'
KLine = require './kline'

colors = ["#000", "#000", "#f00", "#080", "#f00", "#080"]

class KLineSegment
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.segment
    @_ui = @root._ui

  init: ->

  update: (data, datasel, dataset) ->
    @_ui.svg.select("g.segment").remove()
    if not datasel.Segment
      return
    if not datasel.Segment.Data
      return
    g = @_ui.svg.append("g")
      .attr("class", "segment")

    x = @_ui.x
    y = @_ui.y

    dispatch = @root.dispatch
    sdata = datasel.Segment.Data
    dataset = KLine.filter sdata, data
    color = (d, i) -> colors[d.Type] || colors[0]
    g.selectAll('circle')
      .data(dataset)
      .enter()
      .append('circle')
      .attr('cx', (d) -> x d.i)
      .attr('cy', (d) -> y d.Price)
      .attr('r', 6)
      .style("stroke", color)
      .style("fill", (d,i) -> if d.Case1 then color(d,i) else '#fff')
      .on('mouseover', (d, i) -> dispatch.tip @, 'segment', d, i)

KLine.register_plugin 'segment', KLineSegment
