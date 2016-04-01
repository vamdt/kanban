import d3 from 'd3';
import config from '../config';

export default function () {
  const conf = config.load();
  const pool = conf.pool || '';
  const sid = this.$route.params.sid || '';
  d3.json(`/lucky?s=${sid}&pool=${pool}`, (err, data) => {
    if (data && data.lucky) {
      this.$route.router.go({
        name: 'stock',
        params: { sid: data.lucky, k: 1 },
      });
    }
  });
}
