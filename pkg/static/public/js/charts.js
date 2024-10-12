class EChartsLineChart extends HTMLElement {
  constructor() {
    super();
    // Attach shadow DOM
    this.attachShadow({ mode: "open" });
  }

  async connectedCallback() {
    const dataType = this.getAttribute("data-type");
    const range = this.getAttribute("range");
    const deveui = this.getAttribute("deveui");

    const url = `/devices/${deveui}/history?type=${dataType}&range=${range}`;

    let data;
    try {
      // Fetch the data
      const response = await fetch(url);

      // Check if the response is ok
      if (!response.ok) {
        throw new Error("Network response was not ok " + response.statusText);
      }

      // Parse the JSON data
      data = await response.json();

      // Return or process the data
    } catch (error) {
      // Handle errors
      console.error("There was a problem with the fetch operation:", error);
    }

    // Create a container for the chart
    const container = document.createElement("div");
    container.style.width = "700px";
    container.style.height = "450px";
    this.shadowRoot.appendChild(container);

    // Initialize ECharts inside the container
    const chart = echarts.init(container, "dark");

    // Set chart options
    const option = {
      title: {
        left: "center",
        text: dataType,
      },
      tooltip: {
        trigger: "axis",
        position: function (pt) {
          return [pt[0], "15%"];
        },
      },
      toolbox: {
        feature: {
          dataZoom: {
            yAxisIndex: "none",
          },
          restore: {},
          saveAsImage: {},
        },
      },
      dataZoom: [
        {
          type: "inside",
          start: 0,
          end: 100,
        },
        {
          start: 0,
          end: 100,
        },
      ],
      xAxis: {
        type: "time",
      },
      yAxis: {
        type: "value",
        min: function (value) {
          if (dataType == "temp") {
            return Math.min(Math.round(value.min - 2), 15);
          } else if (dataType == "humi") {
            return Math.min(Math.round(value.min - 2), 40);
          }
          return null;
        },
        max: function (value) {
          if (dataType == "temp") {
            return Math.round(value.max + 1);
          } else if (dataType == "humi") {
            return Math.max(Math.round(value.max + 2), 80);
          }
          return null;
        },
        splitNumber: 10,
      },
      series: [
        {
          data: data,
          type: "line",
        },
      ],
    };

    // Set the chart option
    chart.setOption(option);
  }
}

// Define the custom element
customElements.define("echarts-line-chart", EChartsLineChart);
