const defaults = {
  nc: false,
  nmas: false,
  ocl: false,
  nvolume: true,
  nmacd: true,
  mas: [
    { interval: 5 },
    { interval: 10 },
    { interval: 20 },
  ],
  color: {
    up: '#f00',
    down: '#080',
    eq: '#000',
  },
  typing_circle_size: 1,
  segment_circle_size: 3,
};

function load() {
  let s = defaults;
  try {
    s = JSON.parse(localStorage.getItem('settings'));
  } catch (e) {
    s = defaults;
  }
  return s || defaults;
}

function save(settings) {
  return localStorage.setItem('settings', JSON.stringify(settings));
}

function update(settings) {
  const o = load();
  let n = {};
  try {
    n = JSON.parse(JSON.stringify(settings));
  } catch (e) {
    n = {};
  }

  for (const k in n) {
    if (n.hasOwnProperty(k)) {
      o[k] = n[k];
    }
  }
  save(o);
}

export default {
  load,
  save,
  update,
};
