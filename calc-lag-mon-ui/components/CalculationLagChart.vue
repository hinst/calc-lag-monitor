<script lang="ts">

import { defineComponent } from '@nuxtjs/composition-api';
import { Line, mixins } from 'vue-chartjs';

export default defineComponent({
  extends: Line,
  mixins: [mixins.reactiveProp],
  props: {
    options: {
      type: Object
    }
  },
  mounted() {
    (this as any).renderChart(this.chartData, this.options);
  }
})

</script>
