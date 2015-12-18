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

    line = dataset.m1s.Segment.Line
    for level in levels when dataset[level.name]
      k = level.level
      d = level.name
      iscur = ksel is k
      if iscur
        ssel = line
      @draw(k, dataset[d].Hub.Data, line, data, iscur)
      if should_return
        return
      if iscur
        should_return = on
      line = dataset[d].Hub.Line

  draw: (k, data, line, kdata, issel) ->
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
      .attr("width", (d, i) -> Math.max(1, x(d.ei) - x(d.i)))
      .attr("height", (d, i) -> Math.max(1, y(d.Low) - y(d.High)))
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

    if issel and k isnt '1'
      @link(g, line, kdata)

  link: (g, ldata, kdata) ->
    svg = @_ui.svg
    path = g.append("path")
      .style("fill", "none")
      .style("stroke", '#abc')
      .style("stroke-width", "2")

    x = @_ui.x
    y = @_ui.y

    dataset = KLine.filter ldata, kdata
    path.data([dataset])

    line = d3.svg.line()
      .x((d) -> x d.i)
      .y((d) -> y d.Price)

    path.attr("d", line)

KLine.register_plugin 'hub', KLineHub
