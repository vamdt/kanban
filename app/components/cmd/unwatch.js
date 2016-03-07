export default function (...args) {
  if (args.length < 1) {
    return;
  }
  let num = 0;
  args.forEach((o) => {
    const i = this.stocks.findIndex((e) => e.sid === o);
    if (i > -1) {
      this.stocks.splice(i, 1);
      num++;
    }
  });

  if (num) {
    localStorage.setItem('stocks', JSON.stringify(this.stocks));
  }
}
