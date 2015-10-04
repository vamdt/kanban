d3 = require 'd3'
KLine = require './kline'

formatValue = d3.format(",d")
fmtVolume = (d) -> formatValue(d/100) + '手'

class KLineVolume
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.volume

  init: ->
    margin = @root.options.margin
    width = @root.options.width
    height = 50
    @height = height

    rsvg = @root._ui.svg
    container = @root._ui.container

    w = +d3.select(rsvg.node().parentNode).attr("width")
    top = +margin.top+d3.select(rsvg.node().parentNode).attr("height")
    svg = container.append("svg")
      .attr("width", w)
      .attr("height", height)
      .append("g")
      .attr("transform", "translate(#{margin.left},0)")
    @svg = svg

    y = @y = d3.scale.linear()
      .range([height, 0])

    yAxis = @yAxis = d3.svg.axis()
      .scale(y)
      .orient("left")
      .ticks(4)
      .tickSize(-w)
      .tickFormat(fmtVolume)

    svg.append("g")
      .attr("class", "y axis")
      .call(yAxis)

  update: (data) ->
    kColor = KLine.kColor
    x = @root._ui.x
    y = @y
    y.domain([0, d3.max(data, (d)->d.volume)])
    height = @height
    candleWidth = @root.options.candle.width
    svg = @svg

    svg.select(".y.axis").call(@yAxis)
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
