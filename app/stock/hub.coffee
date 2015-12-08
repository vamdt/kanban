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
      @draw(k, dataset[d].Hub.Data, line)
      line = dataset[d].Hub.Line

  draw: (k, data, line) ->
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
    left = @root._left
    size = @root.options.size

    dataset = []
    last = {}
    for d in data
      if not d.oI
        for t in line
          if t.I == d.I
            d.oI = t.oI
          if t.End == d.End
            d.eI = t.oI
      d.i = d.oI - left

      if d.i >= 0 and d.i <= size
        if last.i < 0 or last.i > size
          dataset.push last
        dataset.push d
      else if last.i >= 0 and last.i <= size
        dataset.push d
      last = d

    g.selectAll("rect")
      .data(dataset)
      .enter()
      .append("rect")
      .attr("x", (d, i) -> x(d.i))
      .attr("y", (d, i) -> y(d.High))
      .attr("width", (d, i) -> x(d.eI) - x(d.oI))
      .attr("height", (d, i) -> y(d.Low) - y(d.High))
      .attr("fill", 'steelblue')
      .style("fill-opacity", ".1")

    g.selectAll("text")
      .data(dataset)
      .enter()
      .append("text")
      .attr("x", (d, i) -> x(d.i))
      .attr("y", (d, i) -> y(d.High)+10)
      .attr("fill", 'black')
      .text(k)

KLine.register_plugin 'hub', KLineHub
