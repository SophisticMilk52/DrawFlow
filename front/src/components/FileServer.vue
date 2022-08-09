<template>
  <div class="header">
    <h1>File Server</h1>
    <div class="graph">
      <Apexchart width="500" type="bar" :options="getOptions" :series="getData"></Apexchart>
    </div>
  </div>
</template>
<script>
import Apexchart from "vue3-apexcharts";
import axios from "axios";
export default {
  name: 'FileServer',
  components: {
    Apexchart
  },
  data() {
    return {
      sender: Object,
      subcriber: Object,
      bool: true,
      outputs: [],
      keys: [],
      channels: [],

    };
  }, mounted() {
    // while (this.bool) {
    setInterval(null, 5000)
    this.senderInfo()
    this.subscriberInfo()
    //}
    this.groupBy()
  },
  computed: {
    getOptions() {
      const aux = {
        chart: {
          id: '127.0.0.1:8888',
          type: 'bar',
          height: 350,
          stacked: true,
          toolbar: {
            show: true
          },
          zoom: {
            enabled: true
          }
        },
        xaxis: {
          type: "datetime",
          labels: {
            datetimeUTC: false
          }
        }
      }
      return aux;
    },
    getData() {
      console.log(this.channels)
      console.log(this.outputs)
      var mapa = new Map();
      this.keys.forEach(input => {
        const consolidate = []
        for (let j = 0; j < this.channels.length; j++) {
          var cont = 0
          for (let index = 0; index < this.outputs[input].length; index++) {
     
            if (this.outputs[input][index].channel == this.channels[j]) {
              cont += this.outputs[input][index].size;
            }
            if (index == this.outputs[input].length - 1) {
              consolidate.push([input, cont])
            }

          }
        }
        mapa.set(input, consolidate)
      })
      var aux = []
      var count = 0;
      this.channels.forEach(ch => {
        var arr = []
        this.keys.forEach(k => {
          arr.push(mapa.get(k)[count])
        })
        aux.push({ name: ch, data: arr })
        count++;
      })

      console.log("data", aux)
      return aux
    },
    getSubscriber() {
      return this.subscriber;
    }
  },
  methods: {
    groupBy() {
      const MS_PER_HOURS = 60 * 1000;
      Object.values(this.sender).forEach(input => {
        let time = new Date(input.timestamp).getTime();
        let hour = Math.floor((time) / MS_PER_HOURS);
        const key = MS_PER_HOURS * hour
        if (!this.channels.includes(input.channel)) {
          this.channels.push(input.channel);
        }
        if (!this.outputs[key]) { this.outputs[key] = []; this.keys.push(key); }
        this.outputs[key].push(input);
      });
    },
    senderInfo() {
      axios
        .get(`http://localhost:9090/sender`)
        .then((response) => {
          console.log(response.data)
          this.sender = response.data
          this.groupBy()
        });
    },
    subscriberInfo() {
      axios
        .get(`http://localhost:9090/subscribers`)
        .then((response) => {
          console.log(response.data)
          this.subcriber = response.data
        });
    }
  }


}

</script>
<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-direction: column;
}

.graph {
  display: flex;
  margin-top: 150px;
  flex-direction: row;
  justify-content: space-between;
}

body {
  color: none;
}
</style>