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

  hub: (start, remove) ->
    hub = @datasel.Hub.HCData
    line = @datasel.Segment.HCLine || @datasel.Segment.Line
    if remove
      if start < 0
        return
      while hub.length > start
        hub.pop()
      @save()
      @root.dispatch.redraw()
      return
    begin = @datasel.begin || 0
    begin = +begin
    start = Math.max(start, 1)
    index = start + begin - 1
    if index + 2 >= line.length
      return
    a = line[index]
    b = line[index+1]
    c = line[index+2]
    zg = Math.min(a.High, c.High)
    zd = Math.max(a.Low, c.Low)
    if zd > zg
      return
    h = JSON.parse(JSON.stringify(a))
    h.High = zg
    h.Low = zd
    h.ETime = c.ETime
    for i in ['date', 'edate', 'i', 'ei']
      delete h[i]
    if hub.length > 0
      lh = hub[hub.length-1]
      if lh.Time == h.Time
        hub.pop()
    hub.push h
    @save()
    @root.dispatch.redraw()

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

    prev = off
    for level in levels
      dataset = @dataset[level.name]
      dataset.prev = prev
      prev = dataset

    for level in levels when level.level == @k
      data[level.name] = data[level.name] || {}
      hchub = data[level.name]
      hchub.begin = bnum(hchub)
      hchub.Data = hchub.Data || []
      hchub.Line = hchub.Line || []
      dataset = @dataset[level.name]
      dataset.begin = hchub.begin || 0
      hub = dataset.Hub
      hub.HCData = hub.HCData || hchub.Data
      hub.HCLine = hub.HCLine || hchub.Line

      prev = dataset.prev
      if prev
        dataset.Segment.HCLine = dataset.Segment.HCLine || prev.Hub.HCLine

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
