module.exports = ->
  return if arguments.length < 1
  opts =
    mas: 'nmas'
    candle: 'nc'
    volume: 'nvolume'
    macd: 'nmacd'
    typing: 'ntyping'
    handcraft: 'handcraft'
  v = off
  args = Array.apply(null, arguments)
  if args.length > 1
    if args[args.length-1].toLowerCase() == 'false'
      v = on
      args.pop()
  param = {}
  for o in args
    unless opts[o]
      continue
    param[opts[o]] = v
  @$root.$broadcast('param_change', param)
