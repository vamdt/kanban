css = require 'main.css'
d3 = require 'd3'
KLine = require './kline'
defaults =
  tmpl : '开盘价：<%- p.open%><br/>收盘价：<%- p.close%><br/>最高价：<%- p.high%><br/>最低价：<%- p.low%>'
  margin : [0, 10, 0, 10]
  style:
    'display': 'block'
    'padding': '6px 10px 4px'
    'line-height': '20px'
    'position': 'absolute'
    'font-size': '12px'
    'color': '#fff'
    'background': 'rgb(57, 157, 179)'
    'border-radius': '10px'
    'box-shadow': '0px 0px 3px 3px rgba(0, 0, 0, 0.3)'
    'opacity': '0.8'
    #'visibility': 'hidden'
  width : 120
  height : 82
  x:
    style :
      'display' : 'inline-block'
      'position' : 'absolute'
      'font-size' : '12px'
      'background' : '#FFF'
      'text-align' : 'center'
      'border' : '1px solid #E3F4FF'
    width : 95
    format : 'YYYY/MM/DD HH:mm'
  y:
    style :
      'display' : 'inline-block'
      'position' : 'absolute'
      'font-size' : '12px'
      'background' : '#FFF'
      'border' : '1px solid #E3F4FF'

class KLineToolTip
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.tooltip, defaults

  init: ->
    svg = @root._ui.svg
    container = @root._ui.container
    @root._ui.tips = container.append("div")
      .attr("class", css.tooltip)
      .attr("id", "tooltip")
    for k,v of @options.style
      @root._ui.tips.style k, v


  update: (data) ->
    svg = @root._ui.svg

KLine.register_plugin 'tooltip', KLineToolTip
