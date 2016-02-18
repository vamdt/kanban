d3 = require 'd3'
util = require './util'
defaults =
  container: 'body'
  margin:
    top: 20
    right: 50
    bottom: 30
    left: 50

formatValue = d3.format(",.2f")
fmtCent = (d) -> formatValue d/100

Plugins = {}

class KLine
  constructor: (@options) ->
    @dispatch = d3.dispatch('resize', 'param', 'tip')
    @options = util.extend {}, @options, defaults
    @_data = []
    @_ui = {}
    @_left = 0
    @_max_left = 0
    @plugins = []
    @_param = {}
    @options.size = +@options.size || 100

  update_size: (size, left) ->
    size = size || @options.size || 10
    size = Math.max(size, 10)
    @options.size = size
    left = left || @_left
    atrightedge = @_left == @_max_left
    @_max_left = Math.max(0, @_data.length - 1 - @options.size)
    @options.size = @_data.length - 1 - @_max_left

    if atrightedge and left is @_left
      @_left = @_max_left
    else
      @_left = Math.min(@_max_left, Math.max(0, left))

  data: (data) ->
    @_data = @_data || []
    if not arguments.length
      return @_data.slice(@_left, @_left+@options.size+1)

    return unless data and data.id
    s = data.id
    return if s != @param 's'
    @_dataset = @_dataset || off
    id = @_dataset.id || off
    if id != data.id
      @_dataset = off
    data = util.merge_data(@_dataset, data)
    @_dataset = data
    k = @param 'k'
    dataset = switch k
      when '1' then data.m1s
      when '5' then data.m5s
      when '30' then data.m30s
      when 'week' then data.weeks
      when 'month' then data.months
      else data.days
    @_datasel = dataset
    @_data = dataset.data
    @update_size()

  param: (p) ->
    switch typeof p
      when 'object'
        o = {}
        for k,v of p
          if @_param[k] and @_param[k] != v
            o[k] = @_param[k]
          @_param[k] = v
        if @_param.color
          color[c] = v for c,v of @_param.color
        @dispatch.param(o)
      when 'string'
        return @_param[p]
      else
        return @_param

  init: ->
    @stop()
    @initUI()
    @initPlugins()

    redraw = (data) =>
      return unless data and data.id
      @data data
      @delay_draw()
    @on_event 'kdata', redraw

    @dispatch.on 'param.core', (o) =>
      redraw(@_dataset)

  add_plugin_obj: (plugin) ->
    @plugins.push plugin

  initPlugins: ->
    @plugins = []
    for n,c of Plugins
      plugin = new c @
      continue if plugin.init() is off
      @add_plugin_obj plugin

  resize: (w, h) ->
    w = w || @_ui.container[0][0].clientWidth
    h = h || (util.h() - 80)
    @options.width = w - @options.margin.left - @options.margin.right
    @options.height = h - @options.margin.top - @options.margin.bottom
    @dispatch.resize()

  initUI: ->
    container = @_ui.container = d3.select @options.container || 'body'
    container.html('')
    width = parseInt container.style('width')
    if width < 1
      width = util.w()
    height = parseInt container.style('height')
    if height < 1
      height = util.h() - 80
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
      if data[i].date.getYear() isnt data[prevTick].date.getYear()
        return d3.time.format("%Y-%m-%d")(data[i].date)
      if data[i].date.getMonth() isnt data[prevTick].date.getMonth()
        return d3.time.format("%m-%d")(data[i].date)
      if data[i].date.getDay() isnt data[prevTick].date.getDay()
        return d3.time.format("%d %H:%M")(data[i].date)
      if data[i].date.getHours() isnt data[prevTick].date.getHours()
        return d3.time.format("%H:%M")(data[i].date)
      d3.time.format(":%M")(data[i].date)

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

    zoomed = =>
      n = zoom.scale()
      @zs = @zs || n
      o = @zs
      @zs = n

      x1 = zoom.translate()[0]
      @zx = @zx || x1
      x0 = @zx

      nsize = @options.size
      nleft = @_left
      if n < o
        nsize = parseInt nsize * 1.1
      else if n > o
        nsize = parseInt nsize * 0.9
      else
        return if Math.abs(Math.abs(x1) - Math.abs(x0)) < 2
        @zx = x1

        if x0 > x1
          nleft = nleft + Math.max(20, parseInt nsize * 0.05)
        else if x0 < x1
          nleft = nleft - Math.max(20, parseInt nsize * 0.05)
        else
          return
      @update_size(nsize, nleft)

      @delay_draw()

    zoom = d3.behavior.zoom()
      .on("zoom", zoomed)

    svg
      .call(zoom)
    svg.append("rect")
      .attr("class", "pane")
      .attr("width", width)
      .attr("height", height)

    @dispatch.on 'resize.core', =>
      width = @options.width
      height = @options.height
      margin = @options.margin
      container.select('svg')
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
      x.range([0, width])
      y.range([height, 0])
      xAxis.tickSize(-height, 0)
      yAxis.tickSize(-width)
      svg.select('.x.axis').attr("transform", "translate(0, #{height})")
      svg.select('rect.pane')
        .attr("width", width)
        .attr("height", height)
      @delay_draw()

  move_to: (dir) ->
    switch dir
      when 'left'
        dir = @_left - 100
      when 'right'
        dir = @_left + 100
      when 'home'
        dir = -1
      when 'end'
        dir = 100000
    @update_size(0, +dir)
    @delay_draw()
  delay_draw: ->
    d3.timer => @draw()

  draw: ->
    x = @_ui.x
    y = @_ui.y
    data = @data()

    x.domain([0, data.length-1])
    y.domain([d3.min(data, (d)->d.Low) * 0.99, d3.max(data, (d)->d.High)])

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

    ev =
      s: @param 's'
      k: @param 'k'
      fq: @param 'fq'

    connect = ->
      protocol = if location.protocol.toLowerCase() == 'https:' then 'wss:' else 'ws:'
      ws = new WebSocket("#{protocol}//#{location.host}/socket.io/")
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

    io.on = (event, cb) ->
      handles[event] = handles[event] || []
      i = handles[event].indexOf(cb)
      if i < 0
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
    onclose = (evt) ->
      console.log evt
      setTimeout((->connect()), 1000)

    io.on 'close', onclose

    io.close = ->
      io.off 'close', onclose
      io.ws.close()

    io.connect = connect

  on_event: (event, cb) ->
    @init_websocket =>
      s = @param 's'
      k = @param 'k'
      fq = @param 'fq'
      ename = [s,k,fq,event].join('.')
      @io.on(ename, cb)

  start: ->
    @io.connect()
  stop: ->
    if @io
      @io.close()
      @io = off

  notification: (id, msg) ->
    unless window.Notification
      console.log('bro, your browser is not support notification')
      return
    Notification.requestPermission (permission) ->
      config =
        body: msg
        dir:'auto'
      notification = new Notification(id, config)

KLine.register_plugin = (name, clazz) ->
  Plugins[name] = clazz

color =
  up: "#f00"
  down: "#080"
  eq: "#000"
kColor = (d, i, data) ->
  if d.open == d.close
    if i and data
      if data[i] and data[i-1]
        return color.up if data[i].open >= data[i-1].close
        return color.down if data[i].open < data[i-1].close
    return color.eq
  if d.open > d.close
    return color.down
  color.up

KLine.extend = util.extend
KLine.color = color
KLine.kColor = kColor
KLine.filter = util.filter
KLine.merge_data = util.merge_data
KLine.merge_with_key = util.merge_with_key

module.exports = KLine
