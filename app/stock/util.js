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

export function wwidth() {
  return window.innerWidth || document.documentElement.clientWidth || document.body.clientWidth;
}

export function wheight() {
  return window.innerHeight || document.documentElement.clientHeight || document.body.clientHeight;
}
