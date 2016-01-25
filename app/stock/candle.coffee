d3 = require 'd3'
KLine = require './kline'
KLineMas = require './mas'
defaults =
  width: 2

class KLineCandle
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.candle, defaults
    @root.options.candle = @options
    @_ui = @root._ui

  init: ->
    svg = @root._ui.svg
    container = @root._ui.container
    @options.width = +@options.width || 4
    @root.options.size = Math.floor @root.options.width / (1 + @options.width)
    mas = new KLineMas @root, svg, @root._ui.y, (d) -> d.close
    mas.init()
    @root.add_plugin_obj mas

  update: (data) ->
    kColor = (d, i) -> KLine.kColor d, i, data
    svg = @_ui.svg
    x = @_ui.x
    y = @_ui.y
    candleWidth = @options.width

    nc = @root.param 'nc'
    if nc
      svg.selectAll(".candle").remove()
      return
    dispatch = @root.dispatch
    ocl = @root.param 'ocl'
    if ocl
      svg.selectAll("rect.candle").remove()
    if not ocl
      rect = svg.selectAll("rect.candle")
        .data(data)

      rect
        .enter()
        .append("rect")
        .attr("class", "candle")
        .attr("width", (d, i) -> candleWidth)
        .on('mouseover', (d, i) -> dispatch.tip @, 'k', d, i)

      rect.exit().transition().remove()

      rect
        .transition()
        .attr("x", (d, i) -> x(i) - candleWidth / 2)
        .attr("y", (d, i) -> y(Math.max(d.open, d.close)))
        .attr("height", (d, i) -> Math.max(1, Math.abs(y(d.open) - y(d.close))))
        .attr("stroke", kColor)
        .attr("fill", kColor)

    line = svg.selectAll("line.candle")
      .data(data)

    line
      .enter()
      .append("line")
      .attr("class", "candle")
      .style("stroke-width", "1")
      .on('mouseover', (d, i) -> dispatch.tip @, 'k', d, i)

    line.exit().transition().remove()
    line
      .transition()
      .style("stroke", kColor)
      .attr("x1", (d, i) -> x(i))
      .attr("y1", (d, i) -> y(d.high))
      .attr("x2", (d, i) -> x(i))
      .attr("y2", (d, i) -> y(d.low - Math.max(1, Math.min(d.high - d.low, 0))))

    opacity = @root.param 'opacity'
    if opacity
      svg.selectAll('.candle').transition().style('opacity', opacity)

KLine.register_plugin 'candle', KLineCandle
