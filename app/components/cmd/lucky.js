import d3 from 'd3';

export default function () {
  const sid = this.$route.sid || '';
  d3.json(`/lucky?s=${sid}`, (err, data) => {
    if (data && data.lucky) {
      this.$route.router.go({
        name: 'stock',
        params: { sid: data.lucky, k: 1 },
      });
    }
  });
}
