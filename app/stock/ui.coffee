
class KUI
  constructor: ->
    KLine = require './kline'
    @color = KLine.color

  draw_path: (dataset, id, style) ->
    style = style || {}
    style.fill = style.fill || 'none'
    x = @x
    y = @y
    path = @svg.select("path##{id}")
    if path.empty()
      path = @svg.append("path")
        .attr("id", id)
    path
      .style(style)
      .data([dataset])

    line = d3.svg.line()
      .x((d) -> x d.i)
      .y((d) -> y d.Price)

    path.transition().attr("d", line)

  draw_line: (dataset, clazz, style) ->
    dispatch = @dispatch
    style = style || {}
    style.strokeWidth = style.strokeWidth || '1'
    x = @x
    y = @y
    line = @svg.selectAll("line.#{clazz}")
      .data(dataset)

    line
      .enter()
      .append("line")
      .attr("class", clazz)
      .on("mouseover.tip", (d, i) -> dispatch.tip @, clazz, d, i)
      .style(style)

    up = 4
    down = 5
    yy1 = (d) ->
      y if d.Type == up then d.Low else d.High
    yy2 = (d) ->
      y if d.Type == down then d.Low else d.High
    color = @color
    def_stroke = (d) -> if d.Type == up then color.up else color.down

    line.exit().transition().remove()
    line.transition()
      .attr("x1", (d) -> x d.i)
      .attr("y1", yy1)
      .attr("x2", (d) -> x d.ei)
      .attr("y2", yy2)
      .style("stroke", style.stroke || def_stroke)

  draw_lineno: (dataset, begin, clazz, style) ->
    style = style || {}
    x = @x
    y = @y

    up = 4
    down = 5
    yy2 = (d) ->
      y if d.Type == down then d.Low else d.High
    color = @color
    def_stroke = (d) -> if d.Type == up then color.up else color.down

    data = []
    if dataset.length > 0
      d = dataset[0]
      d = ei: d.i, Type: d.Type, Low: d.Low, High:d.High, no: d.no - 1
      if d.Type == up
        d.Type = down
      else
        d.Type = up
      data = [d]
    dataset.forEach (d) -> data.push(d)
    text = @svg.selectAll("text.#{clazz}")
      .data(data)

    text
      .enter()
      .append("text")
      .attr("class", clazz)
      .style(style)

    numf = (d, i) ->
      n = d.no + 1 - begin
      if n > -1
        n
      else
        ''
    text.exit().transition().remove()
    text.transition()
      .attr("x", (d) -> x d.ei)
      .attr("y", yy2)
      .text(numf)
      .style("stroke", style.stroke || def_stroke)

module.exports = KUI
