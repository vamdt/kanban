import d3 from 'd3';

export default function (sid) {
  if (sid.length < 1) {
    return;
  }

  d3.xhr(`/star?s=${sid}`).send('DELETE');
  this.$root.$broadcast('unwatch', sid);
}
