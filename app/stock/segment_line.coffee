d3 = require 'd3'
KLine = require './kline'

class KLineSegmentLine
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.segment
    @_ui = @root._ui

  init: ->

  update: (data, datasel, dataset) ->
    dataset = KLine.filter datasel.Segment.Line, data
    @root.draw_line(dataset, 'segment_line')

KLine.register_plugin 'segment_line', KLineSegmentLine
