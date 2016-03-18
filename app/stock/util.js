import d3 from 'd3';

export function extend(...args) {
  const dest = {};
  args.forEach((arg) => {
    if (!arg) {
      return;
    }
    Object.keys(arg).forEach((k) => {
      if (!dest.hasOwnProperty(k)) {
        dest[k] = arg[k];
      }
    });
  });
  return dest;
}

const bisect = d3.bisector((d) => +d.date);
const parseDate = d3.time.format('%Y-%m-%dT%XZ').parse;

function indexOfFun(hash, range, start, end) {
  return (date) => {
    const idate = +date;
    if (start > idate) {
      return -1;
    }
    if (idate > end) {
      return hash[end] + 1;
    }
    if (hash.hasOwnProperty(idate)) {
      return hash[idate];
    }
    return bisect.right(range, +date);
  };
}

export function filter(_src, _range) {
  let src = _src;
  if (!Array.isArray(src) || src.length < 1) {
    return [];
  }
  const range = _range;
  if (!Array.isArray(range) || range.length < 2) {
    return [];
  }

  const hash = {};
  range.forEach((d, i) => {
    hash[+d.date] = i;
  });

  src.forEach((d, i) => {
    src[i].no = i;
  });

  const startDate = range[0].date;
  const endDate = range[range.length - 1].date;
  let istart = bisect.left(src, +startDate);
  const iend = bisect.right(src, +endDate);
  istart = Math.max(istart - 1, 0);
  src = src.slice(istart, iend + 1);

  const indexOf = indexOfFun(hash, range, +startDate, +endDate);
  src.forEach((d, i) => {
    src[i].i = indexOf(d.date);
    if (d.ETime) {
      src[i].edate = d.edate || parseDate(d.ETime);
      src[i].ei = indexOf(d.edate);
    }
  });

  return src;
}

export function mergeWithKey(_o, n, k) {
  const o = _o;
  if (!o) {
    return n;
  }
  if (!Array.isArray(n[k]) || n[k].length < 1) {
    return o;
  }
  if (!Array.isArray(o[k]) || o[k].length < 1) {
    o[k] = n[k];
    return o;
  }

  const ndate = +n[k][0].date;
  const odate = +o[k][o[k].length - 1].date;
  const o0date = +o[k][0].date;
  if (odate < ndate) {
    o[k] = o[k].concat(n[k]);
  } else if (o0date >= ndate) {
    o[k] = n[k];
  } else {
    const i = bisect.left(o[k], ndate);
    o[k] = o[k].slice(0, i).concat(n[k]);
  }
  return o;
}

function dataInit(_n) {
  const n = _n;
  if (!n) {
    return n;
  }

  const levels = ['m1s', 'm5s', 'm30s', 'days', 'weeks', 'months'];
  levels.forEach((k) => {
    if (!n[k] || !n[k].data) {
      return;
    }

    n[k].data.forEach((d, i) => {
      n[k].data[i].date = d.date || parseDate(d.Time);
    });

    ['Typing', 'Segment', 'Hub'].forEach((name) => {
      if (!n[k][name]) {
        return;
      }
      ['Data', 'Line'].forEach((dn) => {
        if (!n[k][name][dn]) {
          return;
        }
        n[k][name][dn].forEach((d, i) => {
          n[k][name][dn][i].date = d.date || parseDate(d.Time);
        });
      });
    });
  });
  return n;
}

export function mergeData(_o, _n) {
  let o = _o;
  let n = _n;
  if (!n) {
    return o;
  }
  n = dataInit(n);

  if (!o) {
    return n;
  }
  o = dataInit(o);

  Object.keys(n).forEach((k) => {
    if (typeof n[k] !== 'object') {
      o[k] = n[k];
    }
  });

  const levels = ['m1s', 'm5s', 'm30s', 'days', 'weeks', 'months'];
  levels.forEach((k) => {
    if (!n[k] || !n[k].data) {
      return;
    }
    if (!o[k]) {
      o[k] = n[k];
      return;
    }

    o[k] = mergeWithKey(o[k], n[k], 'data');
    ['Typing', 'Segment', 'Hub'].forEach((name) => {
      if (!n[k][name]) {
        return;
      }
      if (!o[k][name]) {
        o[k][name] = n[k][name];
        return;
      }
      ['Data', 'Line'].forEach((dn) => {
        if (n[k][name][dn]) {
          o[k][name] = mergeWithKey(o[k][name], n[k][name], dn);
        }
      });
    });
  });
  return o;
}

export function wwidth() {
  return window.innerWidth || document.documentElement.clientWidth || document.body.clientWidth;
}

export function wheight() {
  return window.innerHeight || document.documentElement.clientHeight || document.body.clientHeight;
}
