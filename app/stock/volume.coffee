d3 = require 'd3'
KLine = require './kline'
KLineMas = require './mas'

formatValue = d3.format(",d")
fmtVolume = (d) -> formatValue(d/100)

class KLineVolume
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.volume

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
      @svg = container.append("svg")
        .attr("width", w)
        .attr("height", height)
        .append("g")
        .attr("transform", "translate(#{margin.left},0)")
    else
      drag = d3.behavior.drag()
        .on("drag", dragmove)
      @svg = rsvg.append("g")
        .call(drag)

    @y = d3.scale.linear()
      .range([height, 0])

    mas = new KLineMas @root, @svg, @y, (d) -> d.volume
    mas.init()
    @root.add_plugin_obj mas

  hide: ->
    @svg.transition().style('display', 'none')
  show: ->
    @svg.transition().style('display', '')

  updateAxis: (data) ->
    if @root.param 'nvolume'
      @hide()
      return off
    else
      @show()
    svg = @svg
    axis = svg.select('#volume_y_axis')
    if axis.empty()
      axis = svg.append("g")
        .attr("class", "y axis")
        .attr("id", "volume_y_axis")
      @yAxis = d3.svg.axis()
        .scale(@y)
        .orient("left")
        .ticks(4)
        .tickSize(-@w)
        .tickFormat(fmtVolume)
    @y.domain([0, d3.max(data, (d)->d.volume)])
    axis.call(@yAxis)

  update: (data) ->
    if @root.param 'nvolume'
      @hide()
      return off
    else
      @show()

    kColor = (d, i) -> KLine.kColor d, i, data
    x = @root._ui.x
    y = @y
    height = @height
    candleWidth = @root.options.candle.width
    svg = @svg

    rect = svg.selectAll("rect.volume")
      .data(data)

    rect.exit().transition().remove()

    rect
      .enter()
      .append("rect")
      .attr("class", "volume")
      .attr("width", candleWidth)

    rect.transition()
      .attr("x", (d, i) -> x(i) - candleWidth / 2)
      .attr("y", (d, i) -> y(d.volume))
      .attr("height", (d, i) -> height - y(d.volume))
      .attr("stroke", kColor)
      .attr("fill", kColor)

KLine.register_plugin 'volume', KLineVolume
