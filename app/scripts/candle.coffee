defaults =
  width: 4

[cup, cdown, ceq] = ["red", "green", "black"]

kColor = (d) ->
  if d.open == d.close
    return ceq
  if d.open > d.close
    return cdown
  cup

class KLineCandle
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.candle, defaults
    @_ui = @root._ui

  init: ->
    svg = @root._ui.svg
    container = @root._ui.container
    @options.width = +@options.width || 4
    @root.options.size = Math.floor @root.options.width / (3 + @options.width)

  update: (data) ->
    svg = @_ui.svg
    x = @_ui.x
    y = @_ui.y
    candleWidth = @options.width

    tips = @_ui.tips
    show = (d, i) ->
      tips.html("#{d.date}<br/>open: #{d.open}<br/>high: #{d.high}<br/>low: #{d.low}<br/>close: #{d.close}")
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

KLine.register_plugin 'candle', KLineCandle
