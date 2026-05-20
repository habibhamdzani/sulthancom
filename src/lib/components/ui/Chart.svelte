<script lang="ts">
  import ApexCharts from 'apexcharts';

  interface Props {
    options: any;
    series: any[];
    type: string;
    height?: number;
  }

  let { options, series, type, height = 300 }: Props = $props();
  let container: HTMLElement;

  let chart: ApexCharts;

  $effect(() => {
    if (!container) return;

    chart = new ApexCharts(container, {
      ...options,
      series,
      chart: {
        ...options.chart,
        type,
        height,
      },
    });

    chart.render();

    return () => {
      chart.destroy();
    };
  });

  $effect(() => {
    if (chart) {
      chart.updateOptions({ series, ...options });
    }
  });
</script>

<div bind:this={container}></div>
