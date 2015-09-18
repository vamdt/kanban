parseDate = d3.time.format("%Y-%m-%d").parse

class KLineAnnotate
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.annotate
    @_ui = @root._ui
    @_notes = []

  init: ->
    svg = @_ui.svg

    defs = svg.append("defs")
    arrowMarker = defs.append("marker")
      .attr("id", "arrow")
      .attr("markerUnits", "strokeWidth")
      .attr("markerWidth", "12")
      .attr("markerHeight", "12")
      .attr("viewBox", "0 0 12 12")
      .attr("refX", "6")
      .attr("refY", "6")
      .attr("orient", "auto")
    arrow_path = "M2,2 L10,6 L2,10 L6,6 L2,2"
    arrowMarker.append("path")
      .attr("d", arrow_path)
      .attr("fill", "#000")

    @root.on_event 'annotate', (data) =>
      console.log data
      @annotate(data)

  clear_annotate: ->
    @_notes = []
    @_ui.svg.selectAll(".notes").remove()

  annotate: (data) ->
    return if data.code != 200
    s = data.param.s
    notes = data.data[s]
    @clear_annotate()
    for n in notes
      n.x1 = parseDate(n.x1) if n.x1
      n.x2 = parseDate(n.x2) if n.x2

      n.y1 = +n.y1
      n.y2 = +n.y2
      n.id = @_notes.length + 1
      @_notes.push n
    @update()

  update: (data) ->
    svg = @_ui.svg
    if not arguments.length
      data = @root.data()
    for n in @_notes
      switch n.type
        when 'rect' then rv = @draw_annotate_rect(data, n)
        when 'point' then rv = @draw_annotate_point(data, n)
        when 'line' then rv = @draw_annotate_line(data, n)
      if rv == off
        @remove_annotate(n)
    on

  remove_annotate: (n) ->
    @_ui.svg.select("g.notes#note#{n.id}").remove()
  draw_annotate_rect: (data, n) ->
    @_ui.color = @_ui.color || d3.scale.category20()
    i0 = (d3.bisector((d) -> d.date).left)(data, n.x1)
    i1 = (d3.bisector((d) -> d.date).left)(data, n.x2)

    return off if i0 == i1
    return off unless data[i0]
    return off unless data[i1]

    x = @_ui.x
    y = @_ui.y
    x0 = x i0
    width = x(i1) - x0
    hl = []
    if n.y1
      hl.push n.y1
    else
      hl.push data[i0].high, data[i0].low
    if n.y2
      hl.push n.y2
    else
      hl.push data[i1].high, data[i1].low

    high = d3.max(hl)
    low = d3.min(hl)
    y0 = y high
    height = Math.max(1, Math.abs(y(high) - y(low)))

    g = @_ui.svg.select("g.notes#note#{n.id}")
    if g.empty()
      g = @_ui.svg.append("g")
        .attr("class", "notes")
        .attr("id", "note#{n.id}")

    rect = g.select("rect")
    if rect.empty()
      rect = g.append("rect")

    rect
      .attr("x", x0)
      .attr("y", y0)
      .attr("width", width)
      .attr("height", height)
      .style("stroke", @_ui.color n.id)
      .style("fill", "none")

    text = g.select("text")
    if text.empty()
      text = g.append("text")
    text
      .attr("x", x0)
      .attr("y", y0)
      .style("fill", @_ui.color n.id)
      .text(n.comment||'')

  draw_annotate_point: (data, n) ->
    n.x2 = n.x2 || n.x1
    i0 = (d3.bisector((d) -> d.date).left)(data, n.x1)
    i1 = (d3.bisector((d) -> d.date).left)(data, n.x2)

    @_ui.color = @_ui.color || d3.scale.category20()
    return off unless data[i0]
    return off unless data[i1]

    x = @_ui.x
    y = @_ui.y
    x1 = x i0
    y1 = y n.y1
    x2 = x i1
    y2 = y n.y2

    g = @_ui.svg.select("g.notes#note#{n.id}")
    if g.empty()
      g = @_ui.svg.append("g")
        .attr("class", "notes")
        .attr("id", "note#{n.id}")

    circle = g.select("circle")
    if circle.empty()
      circle = g.append("circle")

    circle.attr
      cx: x2
      cy: y2
      r: 5

    line = g.select("line")
    if line.empty()
      line = g.append("line")

    textRotate = if n.y1 == n.y2 then "rotate(90)" else ""

    line
      .attr("x1", x1)
      .attr("y1", y1)
      .attr("x2", x2)
      .attr("y2", y2)
      .attr("stroke", @_ui.color n.id)

    text = g.select("text")
    if text.empty()
      text = g.append("text")
    text
      .attr("transform", "translate(#{x1}, #{y1})#{textRotate}")
      .attr("stroke", @_ui.color n.id)
      .text(n.comment||'')

  draw_annotate_line: (data, n) ->
    i0 = (d3.bisector((d) -> d.date).left)(data, n.x1)
    i1 = (d3.bisector((d) -> d.date).left)(data, n.x2)

    @_ui.color = @_ui.color || d3.scale.category20()
    return off unless data[i0]
    return off unless data[i1]

    x = @_ui.x
    y = @_ui.y
    x0 = x i0
    y0 = y n.y1

    g = @_ui.svg.select("g.notes#note#{n.id}")
    if g.empty()
      g = @_ui.svg.append("g")
        .attr("class", "notes")
        .attr("id", "note#{n.id}")

    line = g.select("line")
    if line.empty()
      line = g.append("line")

    line
      .attr("x1", x i0)
      .attr("y1", y n.y1)
      .attr("x2", x i1)
      .attr("y2", y n.y2)
      .attr("stroke", @_ui.color n.id)
      .attr("marker-end","url(#arrow)")

    text = g.select("text")
    if text.empty()
      text = g.append("text")
    text
      .attr("transform", "translate(#{x0}, #{y0})")
      .attr("stroke", @_ui.color n.id)
      .text(n.comment||'')

KLine.register_plugin 'annotate', KLineAnnotate
