d3 = require 'd3'
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

parseDate = d3.time.format("%Y-%m-%dT%XZ").parse
formatValue = d3.format(",.2f")
fmtCent = (d) -> formatValue d/100

Plugins = {}

[cup, cdown, ceq] = ["#f00", "#080", "#000"]

kColor = (d, i, data) ->
  if d.open == d.close
    if i and data
      if data[i] and data[i-1]
        return cup if data[i].open >= data[i-1].close
        return cdown if data[i].open < data[i-1].close
    return ceq
  if d.open > d.close
    return cdown
  cup

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
      return @_data.slice(@_left, @_left+@options.size+1)

    s = data.id
    k = data.param.k
    dataset = switch k
      when '1' then data.m1s
      when '5' then data.m5s
      when '30' then data.m30s
      when 'week' then data.weeks
      when 'month' then data.months
      else data.days
    @_dataset = data
    @_datasel = dataset
    data = dataset.data

    data.forEach (d) ->
      d.date = parseDate(d.time)
    @_data = data
    @_max_left = Math.max(0, data.length - @options.size)
    @_left = @_max_left if @_left < 1

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
      @data data
      @draw()

  add_plugin_obj: (plugin) ->
    @plugins.push plugin

  initPlugins: ->
    for n,c of Plugins
      plugin = new c @
      plugin.init()
      @add_plugin_obj plugin

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
        return d3.time.format("%H:%M")(data[i].date)
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
      .tickFormat(fmtCent)

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
      if @defer_draw_t
        clearTimeout @defer_draw_t
        @defer_draw_t = undefined

      if @defer_ttl
        @defer_ttl--
      else @defer_ttl = 20
      if @defer_ttl < 1
        @draw()
      else
        defer_draw = => @draw()
        @defer_draw_t = setTimeout defer_draw, 100

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
    y.domain([d3.min(data, (d)->d.low) * 0.99, d3.max(data, (d)->d.high)])

    @update data, @_datasel, @_dataset

  updateAxis: ->
    @_ui.svg.select(".x.axis").call(@_ui.xAxis)
    @_ui.svg.select(".y.axis").call(@_ui.yAxis)

  update: (data, datasel, dataset) ->
    @updateAxis()
    for plugin in @plugins when plugin.updateAxis
      plugin.updateAxis data, datasel, dataset
    for plugin in @plugins
      plugin.update data, datasel, dataset

  init_websocket: (done) ->
    if @io
      @io.on('ready', done)
      if @io.connected
        return done()
      return

    io = {}
    handles = {}
    @io = io

    connect = ->
      ws = new WebSocket("ws://#{location.hostname}:3002/socket.io/")
      io.ws = ws
      ws.onopen = (evt) ->
        io.connected = true
        io.trigger('ready', evt)
        io.ws.send(JSON.stringify(ev))

      ws.onclose = (evt) ->
        io.connected = false
        io.trigger('close', evt)

      ws.onmessage = (evt) ->
        res = JSON.parse(evt.data)
        io.trigger('data', res)

      ws.onerror = (evt) ->
        io.trigger('error', evt)

    ev =
      s: @param 's'
      k: @param 'k'
      fq: @param 'fq'

    io.on = (event, cb) ->
      handles[event] = handles[event] || []
      handles[event].push(cb)

    io.off = (event, cb) ->
      return unless handles[event]
      i = handles[event].indexOf(cb)
      if i > -1
        handles[event].splice i, 1

    io.emit = (event, data) ->
      msg = event: event, data: data
      io.ws.send(JSON.stringify(msg))

    io.trigger = (event, data) ->
      return unless handles[event]
      fn(data) for fn in handles[event]

    io.on 'data', (data) ->
      for event in ['kdata']
        ename = [ev.s,ev.k,ev.fq,event].join('.')
        data.param = ev
        io.trigger ename, data

    if done
      @io.on 'ready', done
    io.on 'close', (evt) ->
      setTimeout((->connect()), 1000)
    connect()

  on_event: (event, cb) ->
    @init_websocket =>
      s = @param 's'
      k = @param 'k'
      fq = @param 'fq'
      ename = [s,k,fq,event].join('.')
      @io.on(ename, cb)

KLine.register_plugin = (name, clazz) ->
  Plugins[name] = clazz

KLine.extend = extend
KLine.kColor = kColor

module.exports = KLine
