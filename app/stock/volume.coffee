require 'main.css'
d3 = require 'd3'
KLine = require './kline'
KLineMas = require './mas'

formatValue = d3.format(",d")
fmtVolume = (d) -> formatValue(d/100)

class KLineVolume
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.volume

  init: ->
    if @root.param 'nvolume'
      return off
    margin = @root.options.margin
    width = @root.options.width
    height = 50
    @height = height

    rsvg = @root._ui.svg
    container = @root._ui.container

    w = +d3.select(rsvg.node().parentNode).attr("width")
    svg = @svg = container.append("svg")
      .attr("width", w)
      .attr("height", height)
      .append("g")
      .attr("transform", "translate(#{margin.left},0)")

    y = @y = d3.scale.linear()
      .range([height, 0])

    mas = new KLineMas @root, svg, y, (d) -> d.volume
    mas.init()
    @root.add_plugin_obj mas

    yAxis = @yAxis = d3.svg.axis()
      .scale(y)
      .orient("left")
      .ticks(4)
      .tickSize(-w)
      .tickFormat(fmtVolume)

    svg.append("g")
      .attr("class", "y axis")
      .call(yAxis)

  updateAxis: (data) ->
    @y.domain([0, d3.max(data, (d)->d.volume)])
    @svg.select(".y.axis").call(@yAxis)

  update: (data) ->
    kColor = (d, i) -> KLine.kColor d, i, data
    x = @root._ui.x
    y = @y
    height = @height
    candleWidth = @root.options.candle.width
    svg = @svg

    svg.selectAll("rect").remove()
    svg.selectAll("rect")
      .data(data)
      .enter()
      .append("rect")
      .attr("x", (d, i) -> x(i) - candleWidth / 2)
      .attr("y", (d, i) -> y(d.volume))
      .attr("width", (d, i) -> candleWidth)
      .attr("height", (d, i) -> height - y(d.volume))
      .attr("stroke", kColor)
      .attr("fill", kColor)

KLine.register_plugin 'volume', KLineVolume
