d3 = require 'd3'
KLine = require './kline'

class KLineHub
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.hub
    @_ui = @root._ui

  init: ->

  update: (data, datasel, dataset) ->
    ksel = @root.param 'k'
    ks =
      '1': 'm1s'
      '5': 'm5s'
      '30': 'm30s'
      'day': 'days'
      'week': 'weeks'
      'month': 'months'

    line = dataset.m1s.Segment.Line
    for k, d of ks when dataset[d]
      @draw(k, dataset[d].Hub.Data, line, data)
      line = dataset[d].Hub.Line

  draw: (k, data, line, kdata) ->
    svg = @_ui.svg
    @_ui.svg.select("g#hub-#{k}").remove()
    if not data
      console.log 'no hub level', k
      return
    if not line
      console.log 'no hub segment or prev line'
      return
    g = @_ui.svg.append("g")
      .attr("id", "hub-#{k}")

    x = @_ui.x
    y = @_ui.y

    dataset = KLine.filter data, kdata
    g.selectAll("rect")
      .data(dataset)
      .enter()
      .append("rect")
      .attr("x", (d, i) -> x(d.i))
      .attr("y", (d, i) -> y(d.High))
      .attr("width", (d, i) -> x(d.ei) - x(d.i))
      .attr("height", (d, i) -> y(d.Low) - y(d.High))
      .attr("fill", 'steelblue')
      .style("stroke", 'green')
      .style("stroke-width", '0')
      .style("fill-opacity", ".1")
      .on("mouseover", -> d3.select(@).style("stroke-width", "1"))
      .on("mouseout", -> d3.select(@).style("stroke-width", "0"))

    g.selectAll("text")
      .data(dataset)
      .enter()
      .append("text")
      .attr("x", (d, i) -> x(d.i))
      .attr("y", (d, i) -> y(d.High)+10)
      .attr("fill", 'black')
      .text(k)

KLine.register_plugin 'hub', KLineHub
