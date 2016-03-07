export default function (...args) {
  if (args.length < 1) {
    return;
  }
  let num = 0;
  args.forEach((o) => {
    const s = {
      sid: o,
      name: o,
    };

    const i = this.stocks.findIndex((e) => e.sid === o);
    if (i > -1) {
      s.name = this.stocks[i].name;
      this.stocks.splice(i, 1);
    }
    this.stocks.unshift(s);
    num++;
  });

  if (num) {
    localStorage.setItem('stocks', JSON.stringify(this.stocks));
  }
}
