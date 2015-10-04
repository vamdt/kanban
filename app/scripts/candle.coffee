d3 = require 'd3'
KLine = require './kline'
defaults =
  width: 4

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
    @root.options.size = Math.floor @root.options.width / (3 + @options.width)

  update: (data) ->
    kColor = KLine.kColor
    svg = @_ui.svg
    x = @_ui.x
    y = @_ui.y
    candleWidth = @options.width

    tips = @_ui.tips
    show = (d, i) ->
      tips.html("#{formatDate(d.date)}<br/>open: #{fmtCent(d.open)}<br/>high: #{fmtCent(d.high)}<br/>low: #{fmtCent(d.low)}<br/>close: #{fmtCent(d.close)}<br/>volume: #{d.volume}")
    svg.selectAll("rect.candle").remove()
    svg.selectAll("line.candle").remove()

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
      .attr("y2", (d, i) -> y(d.low))
    opacity = @root.param 'opacity'
    if opacity
      svg.selectAll('.candle').style('opacity', opacity)

KLine.register_plugin 'candle', KLineCandle
