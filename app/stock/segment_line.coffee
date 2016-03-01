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
    @_ui.draw_line(dataset, 'segment_line')
    if handcraft
      begin = datasel.begin || 0
      @_ui.draw_lineno(dataset, begin, 'segment_line')
    else
      @_ui.svg.selectAll("text.segment_line").remove()

KLine.register_plugin 'segment_line', KLineSegmentLine
