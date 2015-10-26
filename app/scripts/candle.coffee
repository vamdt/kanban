d3 = require 'd3'
KLine = require './kline'
KLineMas = require './mas'
defaults =
  width: 2

formatDate = d3.time.format("%Y-%m-%d %X")
formatValue = d3.format(",.2f")
fmtCent = (d) -> formatValue d/100

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
    nmas = @root.getQuery 'nmas'
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

    tips = @_ui.tips
    show = (d, i) ->
      tips.html("#{formatDate(d.date)}<br/>open: #{fmtCent(d.open)}<br/>high: #{fmtCent(d.high)}<br/>low: #{fmtCent(d.low)}<br/>close: #{fmtCent(d.close)}<br/>volume: #{d.volume}")
    svg.selectAll("rect.candle").remove()
    svg.selectAll("line.candle").remove()

    nc = @root.param 'nc'
    if nc
      return
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
        .on('mouseover', (d, i) -> show d, i)

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
      .on('mouseover', (d, i) -> show d, i)
    opacity = @root.param 'opacity'
    if opacity
      svg.selectAll('.candle').style('opacity', opacity)

KLine.register_plugin 'candle', KLineCandle
