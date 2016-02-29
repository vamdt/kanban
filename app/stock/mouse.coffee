d3 = require 'd3'
KLine = require './kline'

class KLineMouse
  constructor: (@root) ->

  init: ->
    svg = @root._ui.svg
    drag = off
    mousedown = ->
      drag = on
    mouseup = ->
      drag = off
    self = @
    mousemove = ->
      return if drag
      return unless self.data
      m = d3.mouse @

    svg
      .on('mousedown.core', mousedown)
      .on('mouseup.core', mouseup)
      .on('mousemove.core', mousemove)

  update: (data, datasel, dataset) ->
    @data = data
    w = @root.options.candle.width
    x = @root._ui.x

KLine.register_plugin 'mouse', KLineMouse
