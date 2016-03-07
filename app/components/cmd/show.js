export default function (...args) {
  if (args.length < 1) {
    return;
  }
  const opts = {
    mas: 'nmas',
    candle: 'nc',
    volume: 'nvolume',
    macd: 'nmacd',
    typing: 'ntyping',
    handcraft: 'handcraft',
  };

  let v = false;
  if (args[args.length - 1].toLowerCase() === 'false') {
    v = true;
    args.pop();
  }
  if (args[args.length - 1].toLowerCase() === 'true') {
    args.pop();
  }

  const param = {};
  args.forEach((e) => {
    if (!opts[e]) {
      return;
    }
    param[opts[e]] = v;
  });
  this.$root.$broadcast('param_change', param);
}
