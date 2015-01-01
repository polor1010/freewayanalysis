<!DOCTYPE html>
<html lang="en">
  	<head>
  		<meta charset="utf-8">
    	<title>High Way Pridect</title>
    	<link href="./static/css/bootstrap.min.css" rel="stylesheet">
        <link href="./static/css/bootstrap-datetimepicker.min.css" rel="stylesheet">
	</head>
 
<style>
p.two {
    border-style: solid;
    border-color: #ff0000 #000000 #0000ff;
}

.states {
    fill: #aaa;
    stroke: #ff0000;
    stroke-width: 3px;
    p.border-style: solid; 
}

.bar {
  fill: steelblue;
}

.bar:hover {
  fill: brown;
}

.axis {
  font: 10px sans-serif;
}

.axis path,
.axis line {
  fill: none;
  stroke: #000;
  shape-rendering: crispEdges;
}

.x.axis path {
  display: none;
}

</style>

<body>
<script src="http://d3js.org/d3.v3.min.js"></script>
<script src="http://d3js.org/topojson.v1.min.js"></script>
<script type="text/javascript" src="./static/jquery/jquery-1.8.3.min.js" charset="UTF-8"></script>
<script type="text/javascript" src="./static/bootstrap/bootstrap.min.js"></script>
<script type="text/javascript" src="./static/bootstrap/bootstrap-datetimepicker.js" charset="UTF-8"></script>
<script type="text/javascript" src="./static/locales/bootstrap-datetimepicker.fr.js" charset="UTF-8"></script>  

	
<table border=0>
  <tr>
    <td>
      <br>
            <div id="datetimepicker" class="form-group">
                <div class="input-group date form_datetime col-md-10" data-date="1979-09-16T05:25:07Z" data-date-format="dd MM yyyy - HH:ii p" data-link-field="dtp_input1" align="right">
                    <input class="form-control" size="26" type="text" value="" readonly>
                    <span class="input-group-addon"><span class="glyphicon glyphicon-th"></span></span>
                </div>
            </div>
      <svg id="map"></svg></td>
    <td valign = "top">
      <br>
      <svg id="chart"></svg>
      <svg id="chart2"></svg>
      <svg id="chart3"></svg>
      <svg id="chart4"></svg>
      <svg id="chart5"></svg>
      <svg id="chart6"></svg>
    </td>
  </tr>
</table>

<p class="two"></p>

<svg xmlns="http://www.w3.org/2000/svg">
  <defs>
     <linearGradient id = "g1" x1 = "50%" y1 = "50%" x2 = "60%" y2 = "60%">
            <stop stop-color = "green" offset = "0%"/>
            <stop stop-color = "pink" offset = "100%"/>
        </linearGradient>
      </defs>
    <use x = "10" y = "10" xlink:href = "#s1" fill = "url(#g1)"/>
   </svg>

   {{.JsonDatas}}

</body>

<script>
 
var width = 500, height = 650;
var projection = d3.geo.mercator().center([122.479531, 23.998567]).scale(9000);
var path = d3.geo.path().projection(projection);

var svg = d3.select("svg#map").attr("width", width).attr("height", height);

d3.json("./static/twCounty2010.topo.json", function (error, data) {

  svg.selectAll("path")
    .data(topojson.feature(data, data.objects.layer1).features)
    .enter()
    .append("path")
    .attr("d", path)
    .attr("stroke", "black")
    .attr("fill", "#ddc");

  });

 d3.json("./static/geo.json", function(error, data) {  
       
    svg.append("g").selectAll("path")
    .data( data.geometries)
    .enter()
    .append("path")
    .attr("d", path)
    .attr("class", "states");
    //.attr("stroke", "red")
    //.attr("stroke-width","3");

});

drawChart({{.JsonDatas}},"svg#chart");
drawChart({{.JsonDatas}},"svg#chart2");
drawChart({{.JsonDatas}},"svg#chart3");
drawChart({{.JsonDatas}},"svg#chart4");
drawChart({{.JsonDatas}},"svg#chart5");
drawChart({{.JsonDatas}},"svg#chart6");

//drawChart("./static/data.tsv","svg#chart6");

function drawChart(filePath,chartID){

//  d3.tsv(filePath, type, function(error, data) {
  	var data = d3.tsv.parse(filePath);
  	console.log(data.speed);
    width = 360;
    height = 160;
    var margin = {top: 20, right: 30, bottom: 30, left: 40};

    var x = d3.scale.ordinal()
        //.range(0,width);
        .rangeRoundBands([0, width], .3);

    var y = d3.scale.linear()
        .range([height, 0]);

    var xAxis = d3.svg.axis()
        .scale(x)
        .orient("bottom");

    var yAxis = d3.svg.axis()
        .scale(y)//.range([10, 350]);
        .orient("left");

    var svg = d3.select(chartID)
        .attr("width", width + margin.left + margin.right)
        .attr("height", height + margin.top + margin.bottom)
      .append("g")
      .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

      x.domain(data.map(function(d) { return d.time; }));
      y.domain([0, d3.max(data, function(d) { return d.speed; })]);
      console.log(d3.max(data));
      svg.append("g")
          .attr("class", "x axis")
          .attr("transform", "translate(0," + height + ")")
          .call(xAxis)
          .append("text")
          .attr("y",10)
          .attr("x",width-10)
          .text("(時間)")
          ;

      svg.append("g")
          .attr("class", "y axis")
          .call(yAxis)
          .append("text")
          .attr("x",-30)
          .attr("y",-5)
          .text("(時速)");
          //.attr("transform", "rotate(-90)");
          //.attr("y", 16);
          //.attr("dy", ".71em")
          //.style("text-anchor", "end");

      svg.selectAll(".bar")
          .data(data)
          .enter().append("rect")
          .attr("class", "bar")
          .attr("x", function(d) { return x(d.time); })
          .attr("width", x.rangeBand())
          .attr("y", function(d) { return y(d.speed); })
          .attr("height", function(d) { return height - y(d.speed); });

 // }
  //);

  function type(d) {
    d.speed = +d.speed;
    return d;
  }

}

    $('.form_datetime').datetimepicker({
        //language:  'fr',
        weekStart: 1,
        todayBtn:  1,
        autoclose: 1,
        todayHighlight: 1,
        startView: 2,
        forceParse: 0,
        showMeridian: 1
    });
</script> 
</html>
