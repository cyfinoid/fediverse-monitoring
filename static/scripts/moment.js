// import {Chart} from 'https://cdn.jsdelivr.net/npm/chart.js@3.0.0/dist/chart.min.js';
// import ChartDataLabels from 'https://cdn.jsdelivr.net/npm/chartjs-plugin-datalabels@2.0.0';

Chart.register(ChartDataLabels);

export async function softwareChart(url, element) {
  try {
    const response = await fetch(url);
    const data = await response.json();

    const labels = data.map(item => item.name);
    const values = data.map(item => item.name_count);

    createPieChart(labels, values, element);
  } catch (error) {
    console.error("Error fetching data:", error);
  }
}

export async function softwareVersionChart(url, element) {
  try {
    const response = await fetch(url);
    const data = await response.json();

    const labels = data.map(item => (item.name + " " + item.version))
    const values = data.map(item => item.name_count);

    createPieChart(labels, values, element);
  } catch (error) {
    console.error("Error fetching data:", error);
  }
}

export async function tldChart(url, element) {
  try {
    const response = await fetch(url);
    const data = await response.json();

    const labels = data.map(item => item.tld)
    const values = data.map(item => item.tld_count);

    createPieChart(labels, values, element);
  } catch (error) {
    console.error("Error fetching data:", error);
  }
}

export async function nodeCountChart(url, element, type) {
  try {
    const response = await fetch(url);
    const data = await response.json();

    const labels = data.map(item => moment(item.date_time, 'YYYY_MM_DD'))
    const values = data.map(item => item.node_count);

    createLineChart(labels, values, element, type);
  } catch (error) {
    console.error("Error fetching data:", error);
  }
}

export async function dailyCommentChart(url, element, type) {
  try {
    const response = await fetch(url);
    const data = await response.json();

    const labels = data.map(item => moment(item.date_time, 'YYYY_MM_DD'))
    const values = data.map(item => item.total_comments);

    createLineChart(labels, values, element, type);
  } catch (error) {
    console.error("Error fetching data:", error);
  }
}

export async function dailyPostChart(url, element, type) {
  try {
    const response = await fetch(url);
    const data = await response.json();

    const labels = data.map(item => moment(item.date_time, 'YYYY_MM_DD'))
    const values = data.map(item => item.total_posts);

    createLineChart(labels, values, element, type);
  } catch (error) {
    console.error("Error fetching data:", error);
  }
}

export async function dailyUserChart(url, element, type) {
  try {
    const response = await fetch(url);
    const data = await response.json();

    const labels = data.map(item => moment(item.date_time, 'YYYY_MM_DD'))
    const values = data.map(item => item.total_users);

    createLineChart(labels, values, element, type);
  } catch (error) {
    console.error("Error fetching data:", error);
  }
}

export function createPieChart(labels, values, element) {
  var ctx = document.getElementById(element).getContext("2d");
  console.log("Created");
  var myPieChart = new Chart(ctx, {
    type: 'pie',
    data: {
      labels: labels,
      datasets: [{
        data: values,
        backgroundColor: getRandomColors(values.length),
      }]
    },
    // options: {
    //   legend: false,
    //   plugins: {
    //     datalabels: {
    //       display: true,
    //       formatter: (val, ctx) => {
    //         // Grab the label for this value
    //         const label = ctx.chart.data.labels[ctx.dataIndex];

    //         // Format the number with 2 decimal places
    //         const formattedVal = Intl.NumberFormat('en-US', {
    //           minimumFractionDigits: 2,
    //         }).format(val);

    //         // Put them together
    //         return `${label}: ${formattedVal}`;
    //       },
    //       color: '#fff',
    //       backgroundColor: '#404040',
    //     },
    //   },
    // },
    options: {
      plugins: {
        legend: false,
        datalabels: {
          labels: {
            index: {
              color: '#404040',
              backgroundColor: '#fff',
              borderColor: '#fff',
              borderWidth: 2,
              borderRadius: 4,
              padding: 0,
              font: {
                size: 12,
              },
              // See https://chartjs-plugin-datalabels.netlify.app/guide/options.html#option-context
              formatter: (val, ctx) => ctx.chart.data.labels[ctx.dataIndex],
              align: 'top',
            },
            // name: {
            //   color: (ctx) => ctx.dataset.backgroundColor,
            //   font: {
            //     size: 16,
            //   },
            //   backgroundColor: '#404040',
            //   borderRadius: 4,
            //   offset: 0,
            //   padding: 2,
            //   formatter: (val, ctx) => `#${ctx.dataIndex + 1}`,
            //   align: 'end',
            //   anchor: 'end',
            // },
            // value: {
            //   color: '#404040',
            //   backgroundColor: '#fff',
            //   borderColor: '#fff',
            //   borderWidth: 2,
            //   borderRadius: 4,
            //   padding: 0,
            //   align: 'bottom',
            //   anchor: 'end',
            // },
          },
        },
      }
    },

  });

}

export async function createLineChart(labels, counts, element, type) {
  var ctx = document.getElementById(element).getContext("2d");

  // Chart data
  const data = {
    labels: labels,
    datasets: [{
      label: 'Count of ' + type,
      borderColor: 'rgb(192, 37, 37)',
      data: counts,
      fill: false,
    }],
  };

  const config = {
    type: 'line',
    data: data,
    options: {
      scales: {
        x: {
          type: 'time',
          time: {
            unit: 'day',
          },
          title: {
            display: true,
            text: 'Date',
          },
        },
        y: {
          title: {
            display: true,
            text: 'Count of ' + type,
          },
        },
      },
    },
  };
  new Chart(ctx, config);
}
export function getRandomColors(numColors) {

  var colors = [];
  for (var i = 0; i < numColors; i++) {
    colors.push('#' + Math.floor(Math.random() * 16777215).toString(16));
  }
  return colors;
}

