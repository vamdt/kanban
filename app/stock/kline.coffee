d3 = require 'd3'
util = require './util'
Plugin = require('./plugin').default
KUI = require('./ui').default
defaults =
  container: 'body'
  margin:
    top: 20
    right: 50
    bottom: 30
    left: 50

class KLine
  constructor: (@options) ->
    @dispatch = d3.dispatch('resize', 'param', 'tip', 'cmd', 'redraw', 'nameChange')
    @options = util.extend {}, @options, defaults
    @_data = []
    @_ui = new KUI(@)
    @_ui.dispatch = @dispatch
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

  cmd: ->
    @dispatch.cmd.apply @, arguments

  data: (data) ->
    @_data = @_data || []
    if not arguments.length
      return @_data.slice(@_left, @_left+@options.size+1)

    return unless data and data.id
    s = data.id
    return if s != @param 's'
    @_dataset = @_dataset || off
    name = @_dataset.Name || off
    id = @_dataset.id || off
    if id != data.id
      @_dataset = off
    data = util.merge_data(@_dataset, data)
    @_dataset = data
    if data.Name != name
      @dispatch.nameChange(data.Name)
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
        @dispatch.param(o)
      when 'string'
        return @_param[p]
      else
        return @_param

  init: ->
    @stop()
    @_ui.init()
    @initPlugins()

    redraw = (data) =>
      return unless data and data.id
      @data data
      @delay_draw()
    @on_event 'kdata', redraw

    @dispatch.on 'param.core', (o) =>
      redraw(@_dataset)

    @dispatch.on 'redraw.core', =>
      redraw(@_dataset)

  add_plugin_obj: (plugin) ->
    @plugins.push plugin

  initPlugins: ->
    @plugins = []
    Plugin.every (name, c) =>
      plugin = new c @
      return if plugin.init() is off
      @add_plugin_obj plugin

  resize: (w, h) ->
    w = w || @_ui.container[0][0].clientWidth
    h = h || (util.h() - 80)
    @options.width = w - @options.margin.left - @options.margin.right
    @options.height = h - @options.margin.top - @options.margin.bottom
    @dispatch.resize()

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
    data = @data()

    @_ui.update data
    @update data, @_datasel, @_dataset

  update: (data, datasel, dataset) ->
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

module.exports = KLine
