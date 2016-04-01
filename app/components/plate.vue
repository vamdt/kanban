<template>
<div>
<div class="pure-menu pure-menu-horizontal">
    <ul class="pure-menu-list">
        <li v-for="p in plate" class="pure-menu-item pure-menu-has-children pure-menu-allow-hover">
          <a v-link="{ path: '/plate/' + p.Id + '/0'}"
          class="pure-menu-link">{{p.Name}} {{p.Factor}}</a>
          <ul class="pure-menu-children">
              <li v-for="s in p.Sub" class="pure-menu-item">
                <a @click.stop.prevent="show(s)" href="#"
                  class="pure-menu-link">{{s.Name}} {{s.Factor}}</a>
              </li>
          </ul>
        </li>
    </ul>
    <li v-if="title" class="pure-menu-item">
      &gt;&gt;<button class="pure-button" @click="pool(id)">{{title}}</button>
    </li>
</div>
<div class="pure-g">
  <div v-for="i in stocks" class="pure-u-1-5">
    <a v-if="i.Leaf" v-link="{ path: '/s/' + i.Name + '/1'}"
    class="pure-menu-link">{{i.Tag}} {{i.Factor}}</a>
  </div>
</div>
</div>
</template>

<script>
import d3 from 'd3';
import config from './config';
const param = (hash = {}, key) => hash[key];
function noop() {}

export default {
  data() {
    return {
      title: '',
      id: '',
      plate: [],
      stocks: [],
    };
  },

  route: {
    data() {
      const pid = param(this.$route.params, 'pid') || 0;
      const id = param(this.$route.params, 'id') || 0;
      this.rdata(0, () => {
        this.rdata(pid, () => {
          this.rdata(id, (s) => {
            if (s) {
              this.stocks = s;
            }
          });
        });
      });
    },
  },

  methods: {
    pool(id) {
      config.update({ pool: id });
    },

    show(plate) {
      this.title = plate.Name;
      this.id = plate.Id;
      if (!plate.Sub || !plate.Sub.length) {
        this.$route.router.go({
          path: `/plate/${plate.Pid}/${plate.Id}`,
          replace: true,
        });
        return;
      }
      this.stocks = plate.Sub;
    },

    rdata(pids, cb = noop) {
      const pid = +pids;
      if (pid === 0 && this.plate.length) {
        cb();
        return;
      }
      d3.json(`/plate?pid=${pid}`, (error, rdata) => {
        const data = rdata.sort((a, b) => b.Factor - a.Factor);
        data.forEach((d, i) => {
          if (!d.Sub) {
            data[i].Sub = [];
          }
        });
        for (let i = 0; i < this.plate.length; i++) {
          const p = this.plate[i];
          if (+p.Id === pid) {
            p.Sub = data;
            return cb(data);
          }
          if (!p.Sub) {
            continue;
          }
          for (let j = 0; j < p.Sub.length; j++) {
            const s = p.Sub[j];
            if (+s.Id === pid) {
              s.Sub = data;
              return cb(data);
            }
          }
        }
        this.plate = data;
        return cb(data);
      });
    },
  },
};
</script>
