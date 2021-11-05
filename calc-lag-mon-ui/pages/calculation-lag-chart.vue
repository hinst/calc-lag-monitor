<template>
  <div>
    <div>
      <div style="max-width: 160px; display: inline-block">
        <v-menu>
          <template v-slot:activator="{ on, attrs }">
            <v-text-field
              v-model="timeStart"
              label="From"
              prepend-icon="mdi-calendar-arrow-left"
              readonly
              v-bind="attrs"
              v-on="on"
              clearable
              @change="changeTimeRange"
              hide-details="true"
            ></v-text-field>
          </template>
          <v-date-picker
            v-model="timeStart"
            @change="changeTimeRange"
            color="cyan"
          ></v-date-picker>
        </v-menu>
      </div>
      <div style="max-width: 160px; display: inline-block">
        <v-menu>
          <template v-slot:activator="{ on, attrs }">
            <v-text-field
              v-model="timeEnd"
              label="To"
              prepend-icon="mdi-calendar-arrow-right"
              readonly
              v-bind="attrs"
              v-on="on"
              clearable
              @change="changeTimeRange"
              hide-details="true"
            ></v-text-field>
          </template>
          <v-date-picker
            v-model="timeEnd"
            @change="changeTimeRange"
            color="cyan"
          ></v-date-picker>
        </v-menu>
      </div>
    </div>
    <div v-if="chartData" style="max-width: 100%">
      <CalculationLagChart :height="500" :chart-data="chartData" :options="chartOptions" />
    </div>
    <div v-if="aggregationLevelText">
      Current aggregation level: {{aggregationLevelText}} | Count of points: {{ countOfPoints }}
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, useContext } from '@nuxtjs/composition-api';
import lodash from 'lodash';
import { DateTime } from 'luxon';
import { buildParametersString } from '../url';

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

function getAggregationLevelTitle(aggregationLevel: number) {
  switch (aggregationLevel) {
    case 0: return 'none';
    case 1: return 'second';
    case 2: return 'minute';
    case 3: return 'hour';
    case 4: return 'day';
    case 5: return 'month';
    case 6: return 'year';
    default: return 'unknown';
  }
}

export default defineComponent({
  setup() {
    const context = useContext();
    const chartData = ref<any>(null);
    const aggregationLevelText = ref<string | null>(null);
    const countOfPoints = ref<number | null>(null);
    const chartOptions = ref<any>({
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
    });
    const timeStart = ref<string | null>(DateTime.now().minus({days: 7}).toFormat('yyyy-MM-dd'));
    const timeEnd = ref<string | null>(DateTime.now().toFormat('yyyy-MM-dd'));
    const load = async () => {
      const params: { [key: string]: string } = {};
      if (timeStart.value && timeStart.value.length)
        params['start'] = '' + DateTime.fromFormat(timeStart.value, 'yyyy-MM-dd').toMillis();
      if (timeEnd.value && timeEnd.value.length)
        params['end'] = '' + DateTime.fromFormat(timeEnd.value, 'yyyy-MM-dd').toMillis();
      const url = context.env.apiUrl + '/clm/lag' + buildParametersString(params);
      const response = await fetch(url);
      const responseObject = JSON.parse(await response.text());
      aggregationLevelText.value = getAggregationLevelTitle(responseObject.AggregationLevel);
      let responseArray: CalculationLagRow[] = responseObject.Rows;
      countOfPoints.value = responseArray.length;
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
    };
    const changeTimeRange = () => {
      load();
    };
    load();
    setTimeout(load, 60 * 1000);
    return {
      chartData,
      chartOptions,
      aggregationLevelText,
      countOfPoints,
      timeStart,
      timeEnd,
      changeTimeRange,
    };
  },
});
</script>
