d3 = require 'd3'
KLine = require './kline'

class KLineHub
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.hub
    @_ui = @root._ui

  init: ->

  update: (data, datasel, dataset) ->
    ksel = @root.param 'k'
    levels = [
      {level: '1', name: 'm1s'}
      {level: '5', name: 'm5s'}
      {level: '30', name: 'm30s'}
      {level: 'day', name: 'days'}
      {level: 'week', name: 'weeks'}
      {level: 'month', name: 'months'}
    ]

    need = off
    for level,i in levels
      k = level.level
      d = level.name
      hubdata = off
      if dataset[d]
        hubdata = dataset[d].Hub.Data
      if need
        hubdata = off
      @draw(k, hubdata, data)
      need = on if k is ksel

  draw: (k, data, kdata) ->
    x = @_ui.x
    y = @_ui.y

    dispatch = @root.dispatch
    dataset = KLine.filter data, kdata
    rect = @_ui.svg.selectAll("rect.hub-#{k}")
      .data(dataset)
    hover =
      'stroke': 'green'
      'stroke-width': '1'
      'stroke-opacity': '1'
    hout =
      'stroke': 'steelblue'
      'stroke-width': '1'
      'stroke-opacity': '.3'
    rect
      .enter()
      .append("rect")
      .attr("class", "hub-#{k}")
      .attr("fill", 'transparent')
      .style(hout)
      .on("mouseover.stroke", -> d3.select(@).style(hover))
      .on("mouseover.tip", (d, i) -> dispatch.tip @, 'hub', d, i)
      .on("mouseout.stroke", -> d3.select(@).style(hout))

    rect.exit().transition().remove()
    rect.transition()
      .attr("x", (d, i) -> x(d.i))
      .attr("y", (d, i) -> y(d.High))
      .attr("width", (d, i) -> Math.max(1, x(d.ei) - x(d.i)))
      .attr("height", (d, i) -> Math.max(1, y(d.Low) - y(d.High)))

    text = @_ui.svg.selectAll("text.hub-#{k}")
      .data(dataset)
    text
      .enter()
      .append("text")
      .attr("class", "hub-#{k}")
      .attr("fill", 'black')
    text.exit().transition().remove()
    text.transition()
      .attr("x", (d, i) -> x(d.i))
      .attr("y", (d, i) -> y(d.High)+10)
      .text(k)

KLine.register_plugin 'hub', KLineHub
