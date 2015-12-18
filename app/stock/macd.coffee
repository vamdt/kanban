require 'main.css'
d3 = require 'd3'
KLine = require './kline'

bColor = (d) ->
  if d.MACD > 0
    return "#f00"
  return "#080"

formatValue = d3.format(",.3f")
fmtMacd = (d) -> formatValue d/10000

class KLineMacd
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.macd

  init: ->
    if @root.param 'nmacd'
      return off
    margin = @root.options.margin
    width = @root.options.width
    @height = height = 50

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

    yAxis = @yAxis = d3.svg.axis()
      .scale(y)
      .orient("left")
      .ticks(4)
      .tickSize(-w)
      .tickFormat(fmtMacd)

    svg.append("g")
      .attr("class", "y axis")
      .call(yAxis)

    ldiff = svg.select("path.diff.line")
    if ldiff.empty()
      ldiff = svg.append("path")
        .attr("class", "line diff")
    ldiff
      .style("stroke", 'silver')
      .style("stroke-width", "1")

    ldea = svg.select("path.dea.line")
    if ldea.empty()
      ldea = svg.append("path")
        .attr("class", "line dea")
    ldea
      .style("stroke", 'gold')
      .style("stroke-width", "1")

  updateAxis: (data) ->
    min = d3.min data, (d) -> Math.min(d.DIFF, d.DEA, d.MACD)
    max = d3.max data, (d) -> Math.max(d.DIFF, d.DEA, d.MACD)
    @y.domain([min, max])
    @svg.select(".y.axis").call(@yAxis)

  update: (data) ->
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
      .attr("y", (d, i) -> y(Math.max(d.MACD, 0)))
      .attr("width", (d, i) -> candleWidth)
      .attr("height", (d, i) -> Math.abs(y(0) - y(d.MACD)))
      .attr("stroke", bColor)
      .attr("fill", bColor)

    ldiff = svg.select("path.diff.line")
    ldiff.data([data])
    line = d3.svg.line()
      .x((d, i) -> x i)
      .y((d, i) -> y d.DIFF)

    ldiff.attr("d", line)

    ldea = svg.select("path.dea.line")
    ldea.data([data])
    line = d3.svg.line()
      .x((d, i) -> x i)
      .y((d, i) -> y d.DEA)

    ldea.attr("d", line)

KLine.register_plugin 'macd', KLineMacd
