d3 = require 'd3'
util = require './util'
KLine = require './kline'
defaults =
  tmpl : '开盘价：<%- p.open%><br/>收盘价：<%- p.close%><br/>最高价：<%- p.High%><br/>最低价：<%- p.Low%>'
  margin : [0, 10, 0, 10]
  style:
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
    'display': 'none'
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

formatValue = d3.format(",.2f")
fmtCent = (d) -> formatValue d/100

class KLineToolTip
  constructor: (@root) ->
    @options = KLine.extend {}, @root.options.tooltip, defaults

  init: ->
    svg = @root._ui.svg
    container = @root._ui.container
    tips = container.append("div")
      .attr("id", "tooltip")
    for k,v of @options.style
      tips.style k, v
    @tips = tips

    templ = (name, d, i) ->
      switch name
        when 'k'
          "#{d.time}<br/>open: #{fmtCent(d.open)}<br/>high: #{fmtCent(d.High)}<br/>low: #{fmtCent(d.Low)}<br/>close: #{fmtCent(d.close)}<br/>volume: #{d.volume}"
        when 'typing', 'segment'
          "#{d.Time}<br/>high: #{fmtCent(d.High)}<br/>low: #{fmtCent(d.Low)}<br/>#{name}"
        when 'hub'
          "#{d.Time} -- #{d.ETime}<br/>high: #{fmtCent(d.High)}<br/>low: #{fmtCent(d.Low)}<br/>#{name}"
        when 'mas'
          e = d3.select(@)
          "#{e.attr('id')}<div style='background-color:#{e.style('stroke')}'>#{e.attr('id')}</div>#{name}"
        else "no templ"
    @root.dispatch.on 'tip', (e) ->
      args = Array.prototype.slice.call(arguments)
      args.shift() if args.length && args[0] == e
      left = ->
        w = util.w()
        tw = tips[0][0].clientWidth
        if w - tw - d3.event.pageX - 30 < 0
          (d3.event.pageX - tw - 10) + "px"
        else
          (d3.event.pageX) + "px"
      top = ->
        h = util.h()
        th = tips[0][0].clientHeight
        if h - th - d3.event.pageY - 30 < 0
          (d3.event.pageY - th - 30) + "px"
        else
          (d3.event.pageY+30) + "px"
      tips
        .style('display', '')
        .style("left", left)
        .style("top", top)
        .html(templ.apply(e, args))
        .transition()
        .duration(5000)
        .transition()
        .style('display', 'none')

  update: (data) ->

KLine.register_plugin 'tooltip', KLineToolTip
