
function updateNavbar(target) {
    var pathname;

    if(typeof target == "string") {
        pathname = target;
    } else {
        pathname = window.location.hash;
    }

    var selectedTab = $("nav a[href='" + pathname + "']");

    if(!$(selectedTab).hasClass("disabled")) {
        $("nav a.selected").removeClass("selected");
        $(selectedTab).addClass("selected");
    }
}

function setupFilters(updateFunction) {
    $("form button").each(function(index, element) {
        $(element).click(function(e) {
            e.preventDefault();
        })
    });
    var date = new Date();

    $("[name='from']").val(format.dateToStr(util.getPrevMonthStr(date)));
    $("[name='until']").val(format.dateToStr(date));
}

function setupVMenu(tgtId) {
    $("#" + tgtId).find(".v-menu li").each(function(index, element) {
        $(element).click(function(e) {
            $(e.target).siblings().removeClass("selected");
            $(e.target).addClass("selected");
            $("#" + tgtId).find(".v-menu-content").hide();
            var href = $(e.target).attr("href");
            location.href = href;

            $("#" + tgtId).find(".v-menu-content").fadeIn(200);
        })
    });

    $("#" + tgtId).find(".v-menu li:first").click();
}

function setupHighcharts() {
    Highcharts.setOptions({
        lang: {
            weekdays: ["Domingo", "Segunda", "Terça", "Quarta", "Quinta", "Sexta", "Sábado"],
            months: [ "Janeiro" , "Fevereiro" , "Março" , "Abril" , "Maio" , "Junho" , "Julho" , "Agosto" , "Setembro" , "Outubro" , "Novembro" , "Dezembro"],
            shortMonths: [ "Jan" , "Fev" , "Mar" , "Abr" , "Mai" , "Jun" , "Jul" , "Ago" , "Set" , "Out" , "Nov" , "Dez"]
        }
    });
}

function setupChart(chart) {
    $.ajax({
        url: chart.url,
        async: true,
        data: util.getInput(chart.input),
        crossDomain: false,
        success: function(data) {
            var seriesData = format.seriesData(data, chart.id);
            
            var chartOpts = charts.defaultOpts;

            if(chart.type) {
                if(chart.type == "weekly") {
                    var newData = new Array(7);

                    for(var j = newData.length - 1; j >= 0; j--) {
                        newData[j] = null;
                    }

                    chartOpts.plotOptions.series = {};
                    chartOpts.plotOptions.series.pointStart = util.getLastSunday(seriesData[0].data[0][0]);

                    for(var k = seriesData[0].data.length - 1; k >= 0; k--) {
                        newData[new Date(seriesData[0].data[k][0]).getUTCDay()] = seriesData[0].data[k][1];
                    };

                    chartOpts.plotOptions.series.pointInterval = 24 * 3600000;
                    seriesData[0].data = newData;
                }
                if(chart.type == "column") {
                    for(s in seriesData) {
                        seriesData[s].type = "column"
                    }
                }
            }

            chartOpts.chart.title = chart.title;
            chartOpts.chart.renderTo = chart.id;
            //chartOpts.series = seriesData;
            chartOpts.series = [{
                type: "column",
                name: "Visa",
                data: [[1362960000000,16936908],[1363046400000,16639639],[1363132800000,14493456],[1363219200000,10852678],[1363305600000,22998541],[1363564800000,21649870],[1363651200000,18813330],[1363737600000,11023594],[1363824000000,13687630],[1363910400000,14110873],[1364169600000,17897634],[1364256000000,10510489],[1364342400000,11836042],[1364428800000,15820798],[1364774400000,13918927],[1364860800000,18920462],[1364947200000,12971961],[1365033600000,12809234],[1365120000000,13703354],[1365379200000,10752106],[1365465600000,10955445],[1365552000000,13425964],[1365638400000,11736299],[1365724800000,8525897],[1365984000000,11339958],[1366070400000,10920394],[1366156800000,33751923],[1366243200000,23796309],[1366329600000,21759743],[1366588800000,15354206],[1366675200000,23722707],[1366761600000,34636913],[1366848000000,13744165],[1366934400000,27295570],[1367193600000,22868730],[1367280000000,24697766],[1367366400000,18112671],[1367452800000,15072312],[1367539200000,12911684],[1367798400000,17737166],[1367884800000,17276868],[1367971200000,16878434],[1368057600000,14241068],[1368144000000,11958926]]
            },
            {
                type: "column",
                name: "Master Card",
                data: [[1362960000000,16936908],[1363046400000,16639639],[1363132800000,14493456],[1363219200000,10852678],[1363305600000,22998541],[1363564800000,21649870],[1363651200000,18813330],[1363737600000,11023594],[1363824000000,13687630],[1363910400000,14110873],[1364169600000,17897634],[1364256000000,10510489],[1364342400000,11836042],[1364428800000,15820798],[1364774400000,13918927],[1364860800000,18920462],[1364947200000,12971961],[1365033600000,12809234],[1365120000000,13703354],[1365379200000,10752106],[1365465600000,10955445],[1365552000000,13425964],[1365638400000,11736299],[1365724800000,8525897],[1365984000000,11339958],[1366070400000,10920394],[1366156800000,33751923],[1366243200000,23796309],[1366329600000,21759743],[1366588800000,15354206],[1366675200000,23722707],[1366761600000,34636913],[1366848000000,13744165],[1366934400000,27295570],[1367193600000,22868730],[1367280000000,24697766],[1367366400000,18112671],[1367452800000,15072312],[1367539200000,12911684],[1367798400000,17737166],[1367884800000,17276868],[1367971200000,16878434],[1368057600000,14241068],[1368144000000,11958926]]
            },
            {
                type: "column",
                name: "Diners",
                data: [[1362960000000,16936908],[1363046400000,16639639],[1363132800000,14493456],[1363219200000,10852678],[1363305600000,22998541],[1363564800000,21649870],[1363651200000,18813330],[1363737600000,11023594],[1363824000000,13687630],[1363910400000,14110873],[1364169600000,17897634],[1364256000000,10510489],[1364342400000,11836042],[1364428800000,15820798],[1364774400000,13918927],[1364860800000,18920462],[1364947200000,12971961],[1365033600000,12809234],[1365120000000,13703354],[1365379200000,10752106],[1365465600000,10955445],[1365552000000,13425964],[1365638400000,11736299],[1365724800000,8525897],[1365984000000,11339958],[1366070400000,10920394],[1366156800000,33751923],[1366243200000,23796309],[1366329600000,21759743],[1366588800000,15354206],[1366675200000,23722707],[1366761600000,34636913],[1366848000000,13744165],[1366934400000,27295570],[1367193600000,22868730],[1367280000000,24697766],[1367366400000,18112671],[1367452800000,15072312],[1367539200000,12911684],[1367798400000,17737166],[1367884800000,17276868],[1367971200000,16878434],[1368057600000,14241068],[1368144000000,11958926]]
            },
            {
                type: "column",
                name: "American Ex. ",
                data: [[1362960000000,16936908],[1363046400000,16639639],[1363132800000,14493456],[1363219200000,10852678],[1363305600000,22998541],[1363564800000,21649870],[1363651200000,18813330],[1363737600000,11023594],[1363824000000,13687630],[1363910400000,14110873],[1364169600000,17897634],[1364256000000,10510489],[1364342400000,11836042],[1364428800000,15820798],[1364774400000,13918927],[1364860800000,18920462],[1364947200000,12971961],[1365033600000,12809234],[1365120000000,13703354],[1365379200000,10752106],[1365465600000,10955445],[1365552000000,13425964],[1365638400000,11736299],[1365724800000,8525897],[1365984000000,11339958],[1366070400000,10920394],[1366156800000,33751923],[1366243200000,23796309],[1366329600000,21759743],[1366588800000,15354206],[1366675200000,23722707],[1366761600000,34636913],[1366848000000,13744165],[1366934400000,27295570],[1367193600000,22868730],[1367280000000,24697766],[1367366400000,18112671],[1367452800000,15072312],[1367539200000,12911684],[1367798400000,17737166],[1367884800000,17276868],[1367971200000,16878434],[1368057600000,14241068],[1368144000000,11958926]]
            }
            ];

            $("#" + chart.id).highcharts(chartOpts);
            console.log("chart '" + chart.id + "' created.");
        },
        error: function(err) {
            console.log("Error while accessing: " + chart.url);
            console.log(err);
        }
    });
}

function updateCharts() {
    var _tables = [];
    var _table;

    $(".chart").each(function(index, elem) {
        _tables.push($(elem).attr("id"));
    });

    for(var c in _tables) {
        _table = charts[_tables[c]];

        (function(_c) {
            $.ajax({
                url: _c.url,
                async: true,
                data: util.getInput(_c.input),
                crossDomain: false,
                success: function(data) {
                    var chart = $("#" + _c.id).highcharts();
                    var seriesData = format.seriesData(data, _c.id);

                    for(var i = 0; i < seriesData.length; i++) {
                        chart.series[i].setData(seriesData[i].data, false);
                    }

                    chart.redraw();
                    console.log("chart '" + _c.id + "' updated.");
                },
                error: function(err) {
                    console.log("Error while accessing: " + _c.url);
                    console.log(err);
                }
            });
        }(_table));
    }
}

function setupTable(table) {
    $.ajax({
        url: table.url,
        async: true,
        data: util.getInput(table.input),
        crossDomain: false,
        success: function(data) {
            var tableData = format.tableData(data, table.formats, table.order);

            console.log("A \n" + tableData);

            feedTable(table.id, JSON.parse(data).data, table.head, table.extra);
            console.log("table '" + table.id + "' created.");
        },
        error: function(err) {
            console.log("Error while accessing: " + table.url);
            console.log(err);
        }
    });
}

function setupAdvTable(table) {
    var fakeData = {
        data: [
            {
                status: "success",
                date: "15/12/1987 - 21:30", //YYYY-MM-DD HH:MI:SS
                conclusionDate: "17/12/1987 - 10:00", //YYYY-MM-DD HH:MI:SS
                location: "São Paulo - Brasil", //billingAddress.city + " - " + billingAddress.country
                amount: "BRL-999.99", //order.currency + "-"
                pMethod: "Crédito",
                issuer: "Visa"
            },
            {
                status: "cancelled",
                date: "15/12/1987 - 21:30",
                conclusionDate: "17/12/1987 - 10:00",
                location: "São Paulo - Brasil",
                amount: 200.50,
                pMethod: "Boleto",
                issuer: "Bradesco"
            }
        ]
    };

    var tableData = format.tableData(fakeData, table.formats, table.order);

    feedTable(table.id, tableData, table.head, table.extra);
    console.log("table '" + table.id + "' created.");
}

function feedTable(targetId, data, head, extra) {
    var tbody = $("<tbody></tbody>");
    var tr;

    if(head) {
        var _head = $("<thead><tr></tr></thead>");

        for (var i = 0; i < head.length; i++) {
            $(_head).append("<td>" + head[i] + "</td>");
        };

        $("#" + targetId).append(_head);
    }

    for (var i = 0; i < data.length; i++) {
        tr = $("<tr></tr>");

        for (var attr in data[i]) {
            $(tr).append("<td>" + data[i][attr] + "</td>");
        }

        $(tbody).append(tr);
    }

    $("#" + targetId).append(tbody);
}

function updateTables(targetId) {
    var _tables = [];

    if(typeof targetId == "object") {
        _tables.concat(targetId);
    } else if(typeof targetId == "string"){
        _tables.push(targetId);
    } else {
        $("table").each(function(index, elem) {
            _tables.push($(elem).attr("id"));
        });
    }

    for(var t in _tables) {
        (function(_t) {
            console.log(_t);
            $("#" + _t).empty();
            setupTable(tables[_t]);
            console.log("table '" + _t + "' updated.");
        }(_tables[t]));
    }
}

function toggleExpand(target) {
    var tgtSub = $(target).next(".sub");
    var subs = $(target).siblings(".sub");

    if($(tgtSub).is(":visible")) {
        $(tgtSub).hide();
    } else {
        $(subs).hide();
        $(tgtSub).show();
    }
}