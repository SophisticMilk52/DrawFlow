<template>
  <div class="header">
    <h1>File Server</h1>
    <div class="container">

      <div class="graph">
        <h3>Senders</h3>
        <Apexchart width="500" type="bar" :options="getOptions" :series="getData"></Apexchart>
      </div>
      <div class="graph">
        <h3>Receivers</h3>
        <table id="customers">
          <thead>
            <tr>
              <th scope="col">IP</th>
              <th scope="col">Channel</th>
              <th scope="col">Time</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="value in getSubscribers" :key="value.addres">
              <td>{{ value.addres }}</td>
              <td>{{ value.channel }}</td>
              <td>{{ value.timestamp }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>
<script>
import Apexchart from "vue3-apexcharts";
import axios from "axios";

const XAXIS_UPDATE_INTERVAL = 2500
export default {
  name: 'FileServer',
  components: {
    Apexchart
  },
  data() {
    return {
      sender: {},
      subscriber: Object,
      bool: true,
      keys: [],
      channels: [],
      outputs: [],
      currentTimestamp: new Date().getTime(),
      xaxisInterval: null,
      // xAxisTimeWindow in minutes
      xAxisTimeWindow: 15
    };
  },
  mounted() {
    this.setXAxisTimeout()
    this.senderInfo()
    this.subscriberInfo()
    setInterval(() => {
      this.senderInfo()
      this.subscriberInfo()
      this.series()
    }
      , 7 * 1000)
  },
  computed: {
    getSubscribers() {
      return Object.values(this.subscriber);
    },
    getXAxisMin(){
      return new Date(this.currentTimestamp-this.xAxisTimeWindow*60000).getTime()
    },
    getOptions() {
      const aux = {
        chart: {
          id: '127.0.0.1:8888',
          stacked: true,
          horizontal:true,
          toolbar: {
            show: true
          },
          columnWidth: `${Math.round(100/this.xAxisTimeWindow)}%`,
          borderRadius: 10,
          distributed:true,
          zoom: {
            enabled: true
          }
        },
        xaxis: {
          type: "datetime",
          tickAmount: 20,
          min: this.getXAxisMin,
          max: this.currentTimestamp
        }
      }
      return aux;
    },
    getData() {
      return Object.values(this.outputs)
    },
    groupBy() {
      const MS_PER_HOURS = 60 * 1000;
      var outputs = []
      Object.values(this.sender).forEach(input => {
        let time = new Date(input.timestamp).getTime();
        let hour = Math.floor((time) / MS_PER_HOURS);
        const key = MS_PER_HOURS * hour
        if (!this.channels.includes(input.channel)) {
          this.channels.push(input.channel);
        }
        if (!outputs[key]) { outputs[key] = []; this.keys.push(key); }
        outputs[key].push(input);
      });
      return outputs
    },
  },
  unmounted(){
    if(this.xaxisInterval){
      this.xaxisInterval()
    }
  },
  methods: {
    setXAxisTimeout(){
      this.xaxisInterval=setInterval(()=> {
        this.currentTimestamp = new Date().getTime()
      },XAXIS_UPDATE_INTERVAL)
    },
    series() {
      var outputs = this.groupBy
      var mapa = new Map();
      this.keys.forEach(input => {
        const consolidate = []
        for (let j = 0; j < this.channels.length; j++) {
          var cont = 0
          for (let index = 0; index < outputs[input].length; index++) {

            if (outputs[input][index].channel == this.channels[j]) {
              cont += outputs[input][index].size;
            }
            if (index == outputs[input].length - 1) {
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
        aux.push({ name: "canal: " + ch, data: arr })
        count++;
      })
      this.keys = []
      this.outputs=aux
    },
    senderInfo() {
      axios
        .get(`http://localhost:9090/sender`)
        .then((response) => {
          console.log(response.data)
          if (response.data == null) {
            throw Error("null or blanK value in sender API")
          } else {
            this.sender = response.data
          }

        }).catch(error => {
          console.log("Error", error.message)
        });
    },
    subscriberInfo() {
      axios
        .get(`http://localhost:9090/subscribers`)
        .then((response) => {
          console.log(response.data)
          if (response.data == null) {
            throw Error("null or blanK value in subscriber API")
          } else {
            this.subscriber = response.data
          }
        }).catch(error => {
          console.log("Error", error.message)
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

.container {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
}

.graph {
  display: flex;
  margin-top: 150px;
  flex-direction: column;
  justify-content: center;
}

#customers {
  font-family: Arial, Helvetica, sans-serif;
  border-collapse: collapse;
  width: 100%;
}

#customers td,
#customers th {
  border: 1px solid #ddd;
  padding: 8px;
}

#customers tr:nth-child(even) {
  background-color: #f2f2f2;
}

#customers tr:hover {
  background-color: #ddd;
}

#customers th {
  padding-top: 12px;
  padding-bottom: 12px;
  text-align: left;
  background-color: #04AA6D;
  color: white;
}
</style>