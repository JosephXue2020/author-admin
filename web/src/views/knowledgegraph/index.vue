<template>
  <div class="kg-container">
    <div class="kg-text">专家知识图谱</div>
    <div id="main"></div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
// 引入 echarts 核心模块，核心模块提供了 echarts 使用必须要的接口。
import * as echarts from 'echarts';

export default {
  name: 'KnowledgeGraph',
  computed: {
    ...mapGetters([
      'name'
    ])
  },
  created(){
    // 注册必须的组件
    echarts.use([
      TitleComponent,
      TooltipComponent,
      GridComponent,
      DatasetComponent,
      TransformComponent,
      BarChart,
      LabelLayout,
      UniversalTransition,
      CanvasRenderer
    ]);
  },
  mounted(){
    let colors = ['#5470c6','#91cc75','#fac858','#ee6666','#73c0de','#3ba272']
    var chartDom = document.getElementById('main');
    var myChart = echarts.init(chartDom);
    var option;

    option = {
      title: {
          text: '专家评分情况'
        },
        tooltip: {},
        xAxis: {
          data: ['薛钊', '梁燕', '郭丽琴', '赵姐', '吴恩达', '周志华']
        },
        yAxis: {},
        series: [
          {
            name: '评分',
            type: 'bar',
            data: [100, 20, 36, 10, 10, 20],
            itemStyle:{
              color:params => colors[params.dataIndex] || '#fac858'
            }
          }
        ]
    };

    myChart.setOption(option);
  }
}
</script>

<style lang="scss" scoped>
.kg {
  &-container {
    margin: 30px;
  }
  &-text {
    font-size: 30px;
    line-height: 46px;
  }
}
#main{
  width:600px;
  height:400px;
  margin:40px auto;
}
</style>
