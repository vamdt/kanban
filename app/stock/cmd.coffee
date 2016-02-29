d3 = require 'd3'
KLine = require './kline'

cmds = [
  'begin'
  'hub'
]

hd = {}
class KLineCmd
  constructor: (@root) ->
    @root.on 'cmd', => @cmd.apply @, arguments

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

  begin: (sid, date) ->
    date = d3.time.format("%Y-%m-%d %H:%M").parse(date)
    if sid != @dataset.id
      console.log 'begin', sid, '!= dataset.id', @dataset.id
      return
    hd[sid] = hd[sid] || @load(sid)
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
      hchub.Data = hchub.Data || []
      hub = @dataset[level.name]
      hub.begin = date
      hub.HCData = hub.HCData || hchub.Data

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

KLine.register_plugin 'cmd', KLineCmd
