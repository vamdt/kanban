import plugin from './plugin';

const hd = {};

class Cmd {
  constructor(root) {
    this.root = root;
    this.root.dispatch.on('cmd', (...args) => this.cmd.apply(this, args));
  }

  init() {
  }

  hub(start, remove) {
    const hub = this.datasel.Hub.HCData;
    if (remove) {
      if (start < 0) {
        return;
      }
      while (hub.length > start) {
        hub.pop();
      }
      this.save();
      this.root.dispatch.redraw();
      return;
    }

    const begin = this.datasel.begin || 0;
    const index = Math.max(start, 1) + (+begin) - 1;
    const line = this.datasel.Segment.HCLine || this.datasel.Segment.Line;
    if (index + 2 >= line.length) {
      return;
    }
    const a = line[index];
    const c = line[index + 2];
    const zg = Math.min(a.High, c.High);
    const zd = Math.max(a.Low, c.Low);
    if (zd > zg) {
      return;
    }
    const h = JSON.parse(JSON.stringify(a));
    h.High = zg;
    h.Low = zd;
    h.ETime = c.ETime;
    ['date', 'edate', 'i', 'ei'].forEach((i) => delete h[i]);

    if (hub.length > 0) {
      if (hub[hub.length - 1].Time === h.Time) {
        hub.pop();
      }
    }
    hub.push(h);
    this.save();
    this.root.dispatch.redraw();
  }

  save() {
    const sid = this.dataset.id;
    hd[sid] = hd[sid] || {};
    localStorage.setItem(`hc${sid}`, JSON.stringify(hd[sid]));
  }

  load(sid) {
    try {
      return JSON.parse(localStorage.getItem(`hc${sid}`));
    } catch (error) {
      return {};
    }
  }

  initHc(bnum) {
    const sid = this.dataset.id;
    hd[sid] = hd[sid] || this.load(sid) || {};
    const data = hd[sid];
    const levels = [
      {
        level: '1',
        name: 'm1s',
      },
      {
        level: '5',
        name: 'm5s',
      },
      {
        level: '30',
        name: 'm30s',
      },
      {
        level: 'day',
        name: 'days',
      },
      {
        level: 'week',
        name: 'weeks',
      },
      {
        level: 'month',
        name: 'months',
      },
    ];

    let prev = false;
    levels.forEach(({ name }) => {
      const dataset = this.dataset[name];
      dataset.prev = prev;
      prev = dataset;
    });

    levels.forEach(({ level, name }) => {
      if (level !== this.k) {
        return;
      }
      data[name] = data[name] || {};
      const hchub = data[name];
      hchub.begin = bnum(hchub);
      hchub.Data = hchub.Data || [];
      hchub.Line = hchub.Line || [];
      const dataset = this.dataset[name];
      dataset.begin = hchub.begin || 0;
      const hub = dataset.Hub;
      hub.HCData = hub.HCData || hchub.Data;
      hub.HCLine = hub.HCLine || hchub.Line;
      if (dataset.prev) {
        dataset.Segment.HCLine = dataset.Segment.HCLine || dataset.prev.Hub.HCLine;
      }
    });
  }

  begin(bnum) {
    this.initHc(() => bnum);
    this.save();
    this.root.dispatch.redraw();
  }

  cmd(...args) {
    if (args.length < 1) {
      return;
    }
    if (!this.root.param('handcraft')) {
      return;
    }

    const main = args.shift();
    if (!this[main]) {
      return;
    }
    this[main].apply(this, args);
  }

  update(data, datasel, dataset) {
    this.data = data;
    this.datasel = datasel;
    this.dataset = dataset;
    if (!this.root.param('handcraft')) {
      return;
    }
    if (!this.datasel) {
      return;
    }
    this.k = this.root.param('k') || '1';
    this.initHc((hchub) => hchub.begin || 0);
  }
}

plugin.register('cmd', Cmd);
