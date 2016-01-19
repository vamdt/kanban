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

    for level,i in levels
      k = level.level
      d = level.name
      hubdata = off
      if dataset[d]
        hubdata = dataset[d].Hub.Data
      @draw(k, hubdata, data)

  draw: (k, data, kdata) ->
    x = @_ui.x
    y = @_ui.y

    dataset = KLine.filter data, kdata
    rect = @_ui.svg.selectAll("rect.hub-#{k}")
      .data(dataset)
    rect
      .enter()
      .append("rect")
      .attr("class", "hub-#{k}")
      .attr("fill", 'steelblue')
      .style("stroke", 'green')
      .style("stroke-width", '0')
      .style("fill-opacity", ".1")
      .on("mouseover", -> d3.select(@).style("stroke-width", "1"))
      .on("mouseout", -> d3.select(@).style("stroke-width", "0"))

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
