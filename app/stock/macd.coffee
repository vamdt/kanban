d3 = require 'd3'
KLine = require './kline'

bColor = (d) ->
  if d.MACD > 0
    return "#f00"
  return "#080"

formatValue = d3.format(",.3f")
fmtMacd = (d) -> formatValue d/1000

class KLineMacd
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.macd

  init: ->
    @height = height = 50
    rsvg = @root._ui.svg
    w = +d3.select(rsvg.node().parentNode).attr("width")
    @w = w

    dragmove = (d) ->
      d3.select(@)
        .attr("transform", "translate(0, #{d3.event.y})")

    if @options.container
      margin = @root.options.margin
      container = d3.select @options.container
      svg = @svg = container.append("svg")
        .attr("width", w)
        .attr("height", height)
        .append("g")
        .attr("transform", "translate(#{margin.left},0)")
    else
      drag = d3.behavior.drag()
        .on("drag", dragmove)
      svg = @svg = rsvg.append("g")
        .call(drag)

  hide: ->
    @svg.transition().style('display', 'none')
  show: ->
    @svg.transition().style('display', '')

  updateAxis: (data) ->
    if @root.param 'nmacd'
      @hide()
      return off
    else
      @show()
    svg = @svg
    axis = svg.select('#macd_y_axis')
    if axis.empty()
      axis = svg.append("g")
        .attr("class", "y axis")
        .attr("id", "macd_y_axis")
      @y = d3.scale.linear()
        .range([@height, 0])
      @yAxis = d3.svg.axis()
        .scale(@y)
        .orient("left")
        .ticks(4)
        .tickSize(-@w)
        .tickFormat(fmtMacd)

    min = d3.min data, (d) -> Math.min(d.DIFF, d.DEA, d.MACD)
    max = d3.max data, (d) -> Math.max(d.DIFF, d.DEA, d.MACD)
    @y.domain([min, max])

    axis.call(@yAxis)

  update: (data) ->
    if @root.param 'nmacd'
      @hide()
      return off
    else
      @show()
    x = @root._ui.x
    y = @y
    height = @height
    candleWidth = @root.options.candle.width
    svg = @svg

    rect = svg.selectAll("rect.macd")
      .data(data)

    rect.exit().transition().remove()

    rect
      .enter()
      .append("rect")
      .attr("class", "macd")
      .attr("width", candleWidth)

    rect.transition()
      .attr("x", (d, i) -> x(i) - candleWidth / 2)
      .attr("y", (d, i) -> y(Math.max(d.MACD, 0)))
      .attr("height", (d, i) -> Math.abs(y(0) - y(d.MACD)))
      .attr("stroke", bColor)
      .attr("fill", bColor)

    ldiff = svg.select("path#diff")
    if ldiff.empty()
      ldiff = svg.append("path")
        .attr("id", "diff")
        .style("fill", 'none')
        .style("stroke", 'silver')
        .style("stroke-width", "1")
    ldiff
      .data([data])

    line = d3.svg.line()
      .x((d, i) -> x i)
      .y((d, i) -> y d.DIFF)

    ldiff.attr("d", line)

    ldea = svg.select("path#dea")
    if ldea.empty()
      ldea = svg.append("path")
        .attr("id", "dea")
        .style("fill", 'none')
        .style("stroke", 'gold')
        .style("stroke-width", "1")
    ldea
      .data([data])

    line = d3.svg.line()
      .x((d, i) -> x i)
      .y((d, i) -> y d.DEA)

    ldea.attr("d", line)

KLine.register_plugin 'macd', KLineMacd
