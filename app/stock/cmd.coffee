d3 = require 'd3'
KLine = require './kline'

cmds = [
  'begin'
  'hub'
]

hd = {}
class KLineCmd
  constructor: (@root) ->
    @root.dispatch.on 'cmd', => @cmd.apply @, arguments

  init: ->

  hub: (sid, cmd, start, remove) ->
    if cmd != 'set'
      return
    @k

  save: ->
    sid = @dataset.id
    hd[sid] = hd[sid] || {}
    data = hd[sid]
    localStorage.setItem('hc'+sid, JSON.stringify(data))

  load: (sid) ->
    try
      JSON.parse localStorage.getItem 'hc'+sid
    catch
      {}

  init_hc: (bnum) ->
    sid = @dataset.id
    hd[sid] = hd[sid] || @load(sid) || {}
    data = hd[sid]
    levels = [
      {level: '1', name: 'm1s'}
      {level: '5', name: 'm5s'}
      {level: '30', name: 'm30s'}
      {level: 'day', name: 'days'}
      {level: 'week', name: 'weeks'}
      {level: 'month', name: 'months'}
    ]

    for level in levels when level.level == @k
      data[level.name] = data[level.name] || {}
      hchub = data[level.name]
      hchub.begin = bnum(hchub)
      hchub.Data = hchub.Data || []
      dataset = @dataset[level.name]
      dataset.begin = hchub.begin || 0
      hub = dataset.Hub
      hub.HCData = hub.HCData || hchub.Data

  begin: (bnum) ->
    @init_hc -> bnum
    @save()
    @root.dispatch.redraw()

  cmd: ->
    args = Array.apply(null, arguments)
    if args.length < 1
      return
    main = args.shift()
    unless @[main]
      return
    handcraft = @root.param 'handcraft'
    unless handcraft
      return
    @[main].apply @, args

  update: (@data, @datasel, @dataset) ->
    @k = @root.param('k') || '1'
    @init_hc (hchub) -> hchub.begin || 0

KLine.register_plugin 'cmd', KLineCmd
