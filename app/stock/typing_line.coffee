d3 = require 'd3'
KLine = require './kline'

class KLineTypingLine
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.typing
    @_ui = @root._ui

  init: ->

  update: (data, datasel, dataset) ->
    dataset = KLine.filter datasel.Typing.Line, data
    style =
      "stroke-dasharray": "7 7"
      "stroke": "#abc"
    @root.draw_line(dataset, "typing_line", style)

KLine.register_plugin 'typing_line', KLineTypingLine
