d3 = require 'd3'
KLine = require './kline'

class KLineSegmentLine
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.segment
    @_ui = @root._ui

  init: ->

  update: (data, datasel, dataset) ->
    dname = 'Line'
    handcraft = @root.param 'handcraft'
    if handcraft and '1' != @root.param 'k'
      dname = 'HCLine'

    line = datasel.Segment[dname]
    dataset = KLine.filter line, data
    up = 4
    down = 5
    dataset.forEach (d) ->
      return if d.hasOwnProperty('MACD')
      return if d.i < 0
      return if d.ei < 0
      i = d.i
      mup = 0
      mdown = 0
      while i < d.ei and i < data.length
        if data[i].MACD > 0
          mup += data[i].MACD
        else if data[i].MACD < 0
          mdown += data[i].MACD
        i++
      d.MACD = if d.Type == up then mup else mdown

    @_ui.draw_line(dataset, 'segment_line')
    if handcraft
      begin = datasel.begin || 0
      @_ui.draw_lineno(dataset, begin, 'segment_line')
    else
      @_ui.svg.selectAll("text.segment_line").remove()

KLine.register_plugin 'segment_line', KLineSegmentLine
