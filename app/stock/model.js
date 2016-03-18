import EventEmitter from './event';
import { bisector, time } from 'd3';

const bisect = bisector((d) => +d.date);
const parseDate = time.format('%Y-%m-%dT%XZ').parse;

function merge(name, data, Clazz) {
  if (!Array.isArray(data)) {
    return;
  }

  if (data.length < 1) {
    return;
  }

  const o = this[name];
  if (o.length > 0) {
    const ndate = +parseDate(data[0].Time);
    const odate = +o[o.length - 1].date;
    const o0date = +o[0].date;
    if (odate < ndate) {
    } else if (o0date >= ndate) {
      while(this[name].length) {
        this[name].pop();
      }
    } else {
      const i = bisect.left(o, ndate);
      while(this[name].length > i) {
        this[name].pop();
      }
    }
  }

  data.forEach((d) => {
    this[name].push(new Clazz(name, d));
  });
}

class Tdata {
  constructor(name, prop) {
    Object.defineProperty(this, 'name', {
      value: 'k',
    });

    Object.keys(prop).forEach((p) => {
      Object.defineProperty(this, p, {
        enumerable: true,
        value: prop[p],
      });
    });

    Object.defineProperty(this, '_date', {
      writable: true,
    });

    Object.defineProperty(this, 'date', {
      get() {
        this._date = this._date || parseDate(this.Time);
        return this._date;
      },
    });
  }
}

class Typing {
  constructor(name, prop) {
    Object.defineProperty(this, 'name', {
      value: name,
    });

    Object.keys(prop).forEach((p) => {
      Object.defineProperty(this, p, {
        enumerable: true,
        value: prop[p],
      });
    });

    Object.defineProperty(this, '_date', {
      writable: true,
    });

    Object.defineProperty(this, 'date', {
      get() {
        this._date = this._date || parseDate(this.Time);
        return this._date;
      },
    });
  }
}

const tsprops = [
  'Data',
  'Line',
];

class Typings {
  constructor(name) {
    Object.defineProperty(this, 'name', {
      value: name.toLowerCase(),
    });
    tsprops.forEach((prop) => {
      Object.defineProperty(this, prop, {
        enumerable: true,
        value: [],
      });
    });
  }

  assign(data) {
    if (typeof data !== 'object') {
      return;
    }
    tsprops.forEach((name) => {
      merge.call(this, name, data[name], Typing);
    });
  }
}

const lprops = [
  'Typing',
  'Segment',
  'Hub',
];

class Level {
  constructor(name, prev) {
    Object.defineProperty(this, 'name', {
      value: name.toLowerCase(),
    });

    if (prev) {
      Object.defineProperty(this, 'prev', {
        value: prev.toLowerCase(),
      });
    }

    Object.defineProperty(this, 'data', {
      enumerable: true,
      value: [],
    });

    lprops.forEach((prop) => {
      Object.defineProperty(this, prop, {
        enumerable: true,
        value: new Typings(prop),
      });
    });
  }

  assign(data) {
    if (typeof data !== 'object') {
      return;
    }

    merge.call(this, 'data', data.data, Tdata);

    lprops.forEach((name) => {
      this[name].assign(data[name]);
    });
  }
}

const levels = [
  'm1s',
  'm5s',
  'm30s',
  'days',
  'weeks',
  'months',
];

export default class Model extends EventEmitter {
  constructor(id) {
    super();
    Object.defineProperty(this, 'id', {
      value: id,
    });

    Object.defineProperty(this, 'name', {
      value: '',
      writable: true,
    });

    levels.reduce((prev, name) => {
      Object.defineProperty(this, name, {
        enumerable: true,
        value: new Level(name, prev),
      });
      return name;
    }, false);
  }

  assign(data) {
    if (typeof data !== 'object' && this.id !== data.id) {
      return;
    }
    this.name = data.Name;

    levels.forEach((name) => {
      this[name].assign(data[name]);
    });
  }
}
