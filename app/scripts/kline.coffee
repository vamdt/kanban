defaults =
  container: 'body'
  margin:
    top: 20
    right: 50
    bottom: 30
    left: 50

extend = ->
  dest = {}
  for i in arguments when i
    for k,v of i
      dest[k] = dest[k] || v
  dest

parseDate = d3.time.format("%Y-%m-%d").parse
formatValue = d3.format(",.2f")
formatCurrency = (d) -> formatValue(d)

Plugins = {}
Watcher = {}

class KLine
  constructor: (@options) ->
    @options = extend {}, @options, defaults
    @_data = []
    @_ui = {}
    @_left = 0
    @_max_left = 0
    @plugins = []
    @_param = {}
    @options.size = +@options.size || 100

  data: (data) ->
    if not arguments.length
      return @_data.slice(@_left, @_left+@options.size)

    return if data.code != 200
    s = data.param.s
    data = data.data[s]
    data.forEach (d) ->
      d.open = +d.data[0]
      d.close = +d.data[1]
      d.high = +d.data[2]
      d.low = +d.data[3]
      d.date = parseDate(d.time)
    @_data = data
    @_max_left = Math.max(0, data.length - 1 - @options.size)
    @_left = @_max_left

  param: (p) ->
    switch typeof p
      when 'object'
        for k,v of p
          @_param[k] = v
      when 'string'
        return @_param[p]
      else
        return @_param

  init: ->
    @initUI()
    @initPlugins()

    @on_event 'kdata', (data) =>
      console.log data
      @data data
      @draw()

  initPlugins: ->
    for n,c of Plugins
      plugin = new c @
      plugin.init()
      @plugins.push plugin

  initUI: ->
    container = @_ui.container = d3.select @options.container || 'body'
    width = parseInt container.style('width')
    if width < 1
      width = window.screen.availWidth
    height = parseInt container.style('height')
    if height < 1
      height = window.screen.availHeight * 0.618
    @options.width = width - @options.margin.left - @options.margin.right
    @options.height = height - @options.margin.top - @options.margin.bottom
    width = @options.width
    height = @options.height
    margin = @options.margin

    svg = @_ui.svg = container.append("svg")
      .attr("width", width + margin.left + margin.right)
      .attr("height", height + margin.top + margin.bottom)
      .append("g")
      .attr("transform", "translate(#{margin.left},#{margin.top})")

    x = @_ui.x = d3.scale.linear()
      .range([0, width])

    y = @_ui.y = d3.scale.linear()
      .range([height, 0])

    xAxisTickFormat = (i) =>
      data = @data()
      if typeof data[i] == 'undefined'
        return 'F'
      @prevTick = @prevTick || 0
      prevTick = @prevTick
      @prevTick = i

      if i == 0
        return d3.time.format("%Y-%m-%d")(data[i].date)
      if data[i].date == data[prevTick].date
        return ''
      if (+data[i].date) - (+data[prevTick].date) < 86400000
        return d3.time.format("%H-%M")(data[i].date)
      d3.time.format("%Y-%m-%d")(data[i].date)

    xAxis = @_ui.xAxis = d3.svg.axis()
      .scale(x)
      .orient("bottom")
      .tickSize(-height, 0)
      .tickFormat(xAxisTickFormat)

    yAxis = @_ui.yAxis = d3.svg.axis()
      .scale(y)
      .orient("left")
      .ticks(6)
      .tickSize(-width)

    svg.append("g")
      .attr("class", "x axis")
      .attr("transform", "translate(0, #{height})")
      .call(xAxis)

    svg.append("g")
      .attr("class", "y axis")
      .call(yAxis)

    dragmove = =>
      if d3.event.dx == 0
        return
      @_left = Math.min(@_max_left, Math.max(0, @_left - d3.event.dx))
      @draw()

    drag = d3.behavior.drag().on("drag", dragmove)

    svg.append("rect")
      .attr("class", "pane")
      .attr("width", width)
      .attr("height", height)
      .call(drag)

  draw: ->
    x = @_ui.x
    y = @_ui.y
    data = @data()

    x.domain([0, data.length-1])
    y.domain([d3.min(data, (d)->d.low) * 0.8, d3.max(data, (d)->d.high)])

    @update data

  updateAxis: ->
    @_ui.svg.select(".x.axis").call(@_ui.xAxis)
    @_ui.svg.select(".y.axis").call(@_ui.yAxis)

  update: (data) ->
    @updateAxis()
    plugin.update data for plugin in @plugins

  on_event: (event, cb) ->
    s = @param 's'
    k = @param 'k'
    fq = @param 'fq'
    @io = @io || io("http://#{location.hostname}:3002")
    ename = [s,k,fq,event].join('.')
    @io.on(ename, cb)
    @io.emit('watch', ename)

KLine.register_plugin = (name, clazz) ->
  Plugins[name] = clazz

KLine.extend = extend

@KLine = KLine
