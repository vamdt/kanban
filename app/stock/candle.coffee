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
    nmas = @root.param 'nmas'
    if nmas
      return
    mas = new KLineMas @root, svg, @root._ui.y, (d) -> d.close
    mas.init()
    @root.add_plugin_obj mas

  update: (data) ->
    kColor = (d, i) -> KLine.kColor d, i, data
    svg = @_ui.svg
    x = @_ui.x
    y = @_ui.y
    candleWidth = @options.width

    svg.selectAll(".candle").remove()

    nc = @root.param 'nc'
    if nc
      return
    dispatch = @root.dispatch
    ocl = @root.param 'ocl'
    if not ocl
      svg.selectAll("rect.candle")
        .data(data)
        .enter()
        .append("rect")
        .attr("class", "candle")
        .attr("x", (d, i) -> x(i) - candleWidth / 2)
        .attr("y", (d, i) -> y(Math.max(d.open, d.close)))
        .attr("width", (d, i) -> candleWidth)
        .attr("height", (d, i) -> Math.max(1, Math.abs(y(d.open) - y(d.close))))
        .attr("stroke", kColor)
        .attr("fill", kColor)
        .on('mouseover', (d, i) -> dispatch.tip d, i, 'k')

    svg.selectAll("line.candle")
      .data(data)
      .enter()
      .append("line")
      .attr("class", "candle")
      .style("stroke", kColor)
      .style("stroke-width", "1")
      .attr("x1", (d, i) -> x(i))
      .attr("y1", (d, i) -> y(d.high))
      .attr("x2", (d, i) -> x(i))
      .attr("y2", (d, i) -> y(d.low - Math.max(1, Math.min(d.high - d.low, 0))))
      .on('mouseover', (d, i) -> dispatch.tip d, i, 'k')
    opacity = @root.param 'opacity'
    if opacity
      svg.selectAll('.candle').style('opacity', opacity)

KLine.register_plugin 'candle', KLineCandle
