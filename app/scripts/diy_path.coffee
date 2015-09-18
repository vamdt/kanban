parseDate = d3.time.format("%Y-%m-%d").parse

class KLineDiyPath
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.diy_path
    @_ui = @root._ui
    @_diy_line = []

  init: ->
    @root.on_event 'path hl', (data) =>
      console.log data
      @diy_line data

  diy_line: (data) ->
    if not arguments.length
      return @_diy_line

    return if data.code != 200

    name = data.param.name
    s = data.param.s
    data = data.data[name][s]

    if data == null
      @_diy_line = []
      @_ui.svg.selectAll(".diy_line").remove()
      return
    dataset = []
    data.forEach (d) ->
      dataset.push
        x1: parseDate(d[0])
        y1: +d[1]
    @_diy_line = dataset
    @update()

  update: (data) ->
    svg = @_ui.svg
    data = @diy_line()
    g = @_ui.svg.select("g.diy_line")
    if !data or data.length < 1
      if !g.empty()
        g.remove()
      return
    if g.empty()
      g = @_ui.svg.append("g")
        .attr("class", "diy_line")

    path = g.select("path")
    if path.empty()
      path = g.append("path")
    path.style("fill", "none")
      .style("stroke", "blue")
      .style("stroke-width", "1")

    x = @_ui.x
    y = @_ui.y
    left = @root._left
    size = @root.options.size

    dataset = []
    last = {}
    for d in data
      d.oi = d.oi || (d3.bisector((d) -> d.date).left)(@root._data, d.x1)
      d.i = d.oi - left
      if d.i >= 0 and d.i <= size
        if last.i < 0 or last.i > size
          dataset.push last
        dataset.push d
      else if last.i >= 0 and last.i <= size
        dataset.push d
      last = d
    path.data([dataset])

    line = d3.svg.line()
      .x((d) -> x d.i)
      .y((d) -> y d.y1)

    path.attr("d", line)

KLine.register_plugin 'diy_path', KLineDiyPath
