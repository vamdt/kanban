cmds = [
  'sugg'
  'unwatch'
  'watch'
  'show'
  'begin'
  'hub'
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
