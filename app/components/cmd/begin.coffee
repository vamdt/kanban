module.exports = ->
  return if arguments.length < 1
  args = Array.apply(null, arguments)
  args.unshift('begin')
  @$root.$broadcast('kline_cmd', args)
