cmds = [
  'sugg'
  'unwatch'
  'watch'
  'show'
]

kanpan = [
  "begin"
  "hub"
]

handlers = {}
for e in cmds
  handlers[e] = require './'+e

module.exports =
  ready: ->
    vm = @
    for e, func of handlers
      do (e, func) ->
        vm.$on e, (opt) ->
          func.apply vm, opt

    for e in kanpan
      do (e) ->
        vm.$on e, (opt) ->
          opt.unshift e
          vm.$root.$broadcast('kline_cmd', opt)
