<template>
  <div>
    <div v-if="chartData" style="max-width: 100%">
      <CalculationLagChart :height="500" :chart-data="chartData" :options="chartOptions" />
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, useContext } from '@nuxtjs/composition-api';
import * as lodash from 'lodash';

interface CalculationLagItem {
  Min: number;
  Average: number;
  Max: number;
}

interface CalculationLagRow {
  Time: string;
  Cheap: CalculationLagItem;
  Expensive: CalculationLagItem;
}

export default defineComponent({
  setup() {
    const context = useContext();
    const chartData = ref<any>(null);
    const chartOptions = ref<any>(
      {
        responsive: true,
        maintainAspectRatio: false,
        scales: {
          xAxes: [{
            type: 'time',
            time: {
              displayFormats: {
                'millisecond': 'MMM DD',
                'second': 'MMM DD',
                'minute': 'MMM DD',
                'hour': 'MMM DD',
                'day': 'MMM DD',
                'week': 'MMM DD',
                'month': 'MMM DD',
                'quarter': 'MMM DD',
                'year': 'MMM DD',
              }
            },
            gridLines: {
              color: 'rgba(200, 200, 200, 0.5)'
            }
          }],
          yAxes: [{
            gridLines: {
              color: 'rgba(200, 200, 200, 0.5)'
            },
            ticks: {
              min: 0,
              beginAtZero: true
            }
          }]
        }
      }
    );
    setTimeout(async () => {
      const response = await fetch(context.env.apiUrl + '/clm/lag');
      let responseArray: CalculationLagRow[] = JSON.parse(await response.text());
      responseArray = lodash.sortBy(responseArray, row => new Date(row.Time));
      const data = responseArray.map(a => ({ x: new Date(a.Time), y: a.Cheap.Average / 1000_000_000 }));
      chartData.value = {
        labels: responseArray.map(a => new Date(a.Time)),
        datasets: [{
          label: 'average cheap',
          data: data,
          backgroundColor: 'rgba(200, 200, 200, 0.5)'
        }]
      };
    }, 1000);
    return {
      chartData,
      chartOptions
    };
  },
});
</script>
