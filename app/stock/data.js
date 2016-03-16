class EventEmitter {
  constructor() {
    this.events = {};
  }

  on(event, cb) {
    this.events[event] = this.events[event] || [];
    if (this.events[event].indexOf(cb) < 0) {
      this.events[event].push(cb);
    }
  }

  off(event, cb) {
    if (!this.events[event]) {
      return;
    }
    if (!cb) {
      delete this.events[event];
      return;
    }
    const i = this.events[event].indexOf(cb);
    if (i > -1) {
      this.events[event].splice(i, 1);
    }
  }

  emit(event, data) {
    if (!this.events[event]) {
      return;
    }
    this.events[event].forEach((fn) => {
      fn(data);
    });
  }
}

export default class IO extends EventEmitter {
  constructor() {
    super();
    this.auto_reconnect = false;

    this.on('ready', () => {
      this.readySubscribe();
    });
  }

  data(evt) {
    const data = JSON.parse(evt.data);
    if (data && data.id) {
      this.emit(`subscribe.${data.id}`, data);
    }
    this.emit('data', evt);
  }

  connect() {
    if (this.connected) {
      return;
    }
    const protocol = location.protocol.toLowerCase() === 'https:' ? 'wss:' : 'ws:';
    const ws = new WebSocket(`${protocol}//${location.host}/socket.io/`);
    this.ws = ws;
    ws.onopen = (evt) => {
      this.connected = true;
      this.emit('ready', evt);
    };

    ws.onclose = (evt) => {
      this.connected = false;
      this.emit('close', evt);
    };

    ws.onmessage = (evt) => {
      this.data(evt);
    };

    ws.onerror = (evt) => {
      this.emit('error', evt);
    };

    this.auto_reconnect = true;
    this.on('close', () => {
      if (!this.auto_reconnect) {
        return;
      }
      setTimeout(() => {
        this.connect();
      }, 2000);
    });
  }

  close() {
    this.auto_reconnect = false;
    if (this.ws) {
      this.ws.close();
      this.ws = undefined;
    }
  }

  readySubscribe() {
    Object.keys(this.events).forEach((e) => {
      const ee = e.split('.');
      if (ee.length !== 2) {
        return;
      }
      if (ee[0] !== 'subscribe') {
        return;
      }
      this.ws.send(JSON.stringify({ s: ee[1] }));
    });
  }

  subscribe(sid, cb) {
    this.on(`subscribe.${sid}`, cb);
    if (this.connected) {
      this.ws.send(JSON.stringify({ s: sid }));
      return;
    }

    if (!this.auto_reconnect) {
      this.connect();
    }
  }

  unsubscribe(sid) {
    this.off(`subscribe.${sid}`);
  }
}
