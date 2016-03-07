import d3 from 'd3';

export default function (sid) {
  if (sid.length < 1) {
    return;
  }
  const s = this.stocks.find((e) => e.sid === sid);
  if (s && s.sid === sid) {
    this.show_stock(s);
    return;
  }

  d3.text(`/search?s=${sid}`, (error, data) => {
    if (error) {
      return;
    }
    const info = data.split(';');
    info.forEach((v, i) => {
      const vv = v.split(',');
      info[i] = {
        sid: vv[3],
        name: vv[4],
      };
    });

    if (info.length === 1) {
      this.show_stock(info[0]);
      return;
    }
    this.sugg = info;
  });
}
