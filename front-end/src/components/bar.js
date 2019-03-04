import {Bar, mixins} from 'vue-chartjs'
const {reactiveProp} = mixins;
export default {
    extends: Bar,
    mixins: [reactiveProp],
    props: ['options'],
    mounted: function () {
        console.log("show graph");
        this.renderChart(this.chartData,this.options);
    }
}
