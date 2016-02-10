'use strict';


techPayControllers.controller('ReportQuantityCtrl',  function($scope, $location, $cookieStore, report) {
    var seriesQuantity = [];
    var chart;
    var seriesVisa;
    var seriesMastercard;
    var seriesDiners;
    var seriesAmex;
    var authToken;

    $scope.main = false;
    authToken = $cookieStore.get('authToken');
    createMonitorChart($location,$scope, report);

    Highcharts.setOptions(Highcharts.theme);

    function load(){
        var divLoading = angular.element(document.querySelector('#loadingContainer'));

        var success = function (data) {
            divLoading.remove(); // remove a div de animacao do DOM
            $scope.main = true;

            switch (data.status) {
                case 'success':

                    var formattedData = format.seriesData( data );

                    console.log( JSON.stringify(formattedData))
                    seriesVisa.setData( findFormattedData( 'VISA', formattedData ),false);
                    seriesMastercard.setData( findFormattedData( 'MASTERCARD', formattedData ),false);
                    seriesDiners.setData( findFormattedData( 'DINERS', formattedData ),false);
                    seriesAmex.setData( findFormattedData( 'AMEX', formattedData ),true);

                    break;
                default:
                    if (data.errorCode == 5020) {
                        $location.path('/start');
                    }
                    alert(data.errorMessage);
                    break;
            }
        };


        var error = function () {
            // TODO: apply user notification here..
        };
        report.quantity(authToken, '2014-01-03','2015-01-05').success(success).error(error);

    }
    function findFormattedData( name, values ){
        for( var i=0;i<values.length;i++){
            if( values[i].name==name ){

                values[i].data.sort(function(x, y){
                    if (x[0] < y[0]) {
                        return -1;
                    }
                    if (x[0] > y[0]) {
                        return 1;
                    }
                    return 0;
                });

                return values[i].data;
            }
        }
        return [];
    }
    function createMonitorChart() {
        $(document).ready(function () {
            Highcharts.setOptions({
                global: {
                    useUTC: false
                }
            });

            chart  = $('#transactionQuantity').highcharts('StockChart', {
                chart: {
                    animation: Highcharts.svg, // don't animate in old IE
                    marginRight: 10,
                    events: {
                        load: function () {
                            seriesVisa = this.series[0];
                            seriesMastercard = this.series[1];
                            seriesDiners = this.series[2];
                            seriesAmex = this.series[3];

                            load();
                        }
                    }
                },
                lang: {
                    weekdays: ["Domingo", "Segunda", "Terça", "Quarta", "Quinta", "Sexta", "Sábado"],
                    months: [ "Janeiro" , "Fevereiro" , "Março" , "Abril" , "Maio" , "Junho" , "Julho" , "Agosto" , "Setembro" , "Outubro" , "Novembro" , "Dezembro"],
                    shortMonths: [ "Jan" , "Fev" , "Mar" , "Abr" , "Mai" , "Jun" , "Jul" , "Ago" , "Set" , "Out" , "Nov" , "Dez"]
                },

                title: {
                    text: 'Quantidade de transações aprovadas'
                },
                xAxis: {
                    type: 'datetime',
                    tickPixelInterval: 150
                },
                yAxis: {
                    title: {
                        text: 'Valor'
                    },
                    plotLines: [{
                        value: 0,
                        width: 1,
                        color: '#780000'
                    }]
                },
                tooltip: {
                    formatter: function () {
                        return '<b>' + seriesQuantity.name + '</b><br/>' +
                            Highcharts.dateFormat('%Y-%m-%d %H:%M:%S', this.x) + '<br/>' +
                            this.y;
                    }
                },
                legend: {
                    enabled: true
                },
                exporting: {
                    enabled: true
                },
                credits: {
                    enabled: false

                },
                series: [
                    {
                        type: "column",
                        name: "Visa"
                    },
                    {
                        type: "column",
                        name: "MasterCard"
                    },
                    {
                        type: "column",
                        name: "Diners"
                    },
                    {
                        type: "column",
                        name: "American Ex. "
                    }
                ]
            });
        });
    }



});

