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
    first = 0
    if datasel.date
      for d,i in line
        if d.date == datasel.date
          first = i
          break
      i = first
      num = 1
      while i < line.length
        line[i].no = num
        i++
        num++
    dataset = KLine.filter line, data
    @root.draw_line(dataset, 'segment_line')
    @root.draw_lineno(dataset, 'segment_line')

KLine.register_plugin 'segment_line', KLineSegmentLine
