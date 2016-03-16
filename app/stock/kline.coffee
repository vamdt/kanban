d3 = require 'd3'
util = require './util'
Plugin = require('./plugin').default
KUI = require('./ui').default
Data = require('./data').default
defaults =
  container: 'body'
  margin:
    top: 20
    right: 50
    bottom: 30
    left: 50

class KLine
  constructor: (@options) ->
    @dispatch = d3.dispatch('resize', 'param', 'tip', 'cmd', 'redraw', 'nameChange', 'uiInit')
    @options = util.extend {}, @options, defaults
    @_data = []
    @_ui = new KUI(@)
    @_ui.dispatch = @dispatch
    @io = new Data()
    @_left = 0
    @_max_left = 0
    @plugins = []
    @_param = {}
    @options.size = +@options.size || 100

    @bindEvent()

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
          if @_param[k] != v
            o[k] = @_param[k]
          @_param[k] = v
        @dispatch.param(o)
      when 'string'
        return @_param[p]
      else
        return @_param

  bindEvent: ->
    redraw = (data) =>
      return unless data and data.id
      @data data
      @delay_draw()

    @dispatch.on 'param.core', (o) =>
      if o.hasOwnProperty('s')
        @io.subscribe @param('s'), redraw
      redraw(@_dataset)

    @dispatch.on 'redraw.core', =>
      redraw(@_dataset)

    @dispatch.on 'uiInit.core', =>
      @initPlugins()

  add_plugin_obj: (plugin) ->
    @plugins.push plugin

  initPlugins: ->
    @plugins = []
    Plugin.every (name, c) =>
      plugin = new c @
      return if plugin.init() is off
      @add_plugin_obj plugin

  resize: (w, h) ->
    unless @_ui.__inited
      return

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

  stop: ->
    @io.close()

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
