var seriesQuantity = [];
var seriesVolume = [];
volumeItems = [];
volumeQtde = [];

techPayControllers.controller('MonitorCtrl',  function($scope, $location, $cookieStore, $interval, report) {
    var divLoading = angular.element( document.querySelector( '#loadingContainer' ) );

    $scope.main = false;

    Highcharts.setOptions(Highcharts.theme);

    createMonitorChart($scope.transactionVolume);

    fillblanc();

    var authToken = $cookieStore.get('authToken');

    var success = function (data) {
        divLoading.remove(); // remove a div de animacao do DOM
        $scope.main = true;

        switch (data.status) {
            case 'success':
                for (var i = 0; i < data.data.length; i++) {
                    var transaction = data.data[i];

                    var x = Math.floor((transaction.gwtime) / 60000 * 1000) * 60000
                    var y = (transaction.amount / 100.00);

                    console.log(x);
                    addVolumeItem(x, y);
                    addVolumeQtde(x,1);

		    var installments="";
		    if(  transaction.request_json.installments ){
			installments =  transaction.request_json.installments.quantity;
		    }
		    var truncated_number="";
		    var brand="";
		    if( transaction.request_json.creditCard ){
			truncated_number = transaction.request_json.creditCard.truncated_number;
			brand = transaction.request_json.creditCard.brand;
		    }

                    $scope.items.push({
                        date: transaction.gwtime * 1000,
                        purchaseCode: transaction.request_json.order.purchaseCode,
                        status: transaction.status,
                        amount: transaction.amount,
                        installments: installments,
                        truncated_number: truncated_number,
                        brand: brand,
                        method: transaction.request_json.order.method,
                        transactionId: transaction.transactionId
                    });
                }
                ;
                //volumeItems.sort(function(x, y){
                //    if (x[0] < y[0]) {
                //        return -1;
                //    }
                //    if (x[0] > y[0]) {
                //        return 1;
                //    }
                //    return 0;
                //});
                //
                //volumeQtde.sort(function(x, y){
                //    if (x[0] < y[0]) {
                //        return -1;
                //    }
                //    if (x[0] > y[0]) {
                //        return 1;
                //    }
                //    return 0;
                //});
                volumeQtde.sort();
                volumeItems.sort();

                seriesVolume.setData(JSON.parse(JSON.stringify(volumeItems)));
                seriesQuantity.setData(JSON.parse(JSON.stringify(volumeQtde)));

                $interval( fillblanc, 60000 );

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




//websocket

    report.listTransaction(authToken,"","",100).success(success).error(error);

    var client = new Faye.Client('https://websocket.techpay.com.br/1.0/ws_transactions', {
        timeout: 120, retry: 5
    });


    var authToken  = $cookieStore.get('authToken');
    var merchantId = $cookieStore.get('merchantId');
    $scope.items=[];


    var clientAuth = {
        outgoing: function(message, callback) {
            if (message.channel !== '/meta/subscribe')
                return callback(message);

            if (!message.ext) message.ext = {};

            message.ext.authToken = authToken;

            callback(message);
        }
    };
    client.addExtension(clientAuth);

    var subscription = client.subscribe('/transactions_' + merchantId, function(message) {
        var div = document.getElementById('txs');

        var transaction = JSON.parse( message.text );

        var status ="Desconhecido";
        switch( transaction.result.status ){
            case  'success':
                status = "Aprovada";
                break;
            case 'error_payment':
                status = "Recusada";
                break;
            case 'pending':
                status = "Processando";
                break;
        }

        var x = Math.floor((transaction.result.time)/60000)*60000
        var y = (transaction.amount/100.00);
        console.log( JSON.stringify(transaction))

        if( transaction.result.status=='success') {
            addVolumeItem(x, y);
            addVolumeQtde(x, 1);
            seriesVolume.setData(JSON.parse(JSON.stringify(volumeItems)));
            seriesQuantity.setData(JSON.parse(JSON.stringify(volumeQtde)));
        }
        var simplifyedJson = {
            date: transaction.gwtime,
            purchaseCode: transaction.request.order.purchaseCode,
            status: transaction.result.status,
            amount: transaction.amount,
            method: transaction.request.order.method,
            transactionId: transaction.transactionId
        }
        if (transaction.request.creditCard && transaction.request.creditCard.truncated_number) {
            simplifyedJson.truncated_number = transaction.request.creditCard.truncated_number;
        }
        if (transaction.request.installments && transaction.request.installments.quantity) {
            simplifyedJson.installments = transaction.request.installments.quantity;
        }
        if (transaction.request.creditCard && transaction.request.creditCard.brand) {
            simplifyedJson.brand = transaction.request.creditCard.brand;
        }
        if(! changeTableItem( $scope.items,simplifyedJson.transactionId, transaction.result.status,transaction.amount )){
            $scope.items.unshift( simplifyedJson );
        }
        $scope.$apply();
    });

});
function changeTableItem(items, transactionId, status, amount)
{
    var found=false;
    for( var i=0;i<items.length;i++ ){
        if( items[i].transactionId == transactionId ){
            items[i].status = status;
            items[i].amount = amount;
            found = true;
        }
    }
    return found;
}

function fillblanc()
{
    console.log("loop");

    var now = new Date().getTime();

    for( var i=0;i<100;i++){
        var x = Math.floor((now-(i*60000)) / 60000 ) * 60000


        console.log(x);

        addVolumeQtde( x, 0);
        addVolumeItem( x, 0 );
    }

    volumeItems.sort(function(x, y){
        if (x[0] < y[0]) {
            return -1;
        }
        if (x[0] > y[0]) {
            return 1;
        }
        return 0;
    });

    volumeQtde.sort(function(x, y){
        if (x[0] < y[0]) {
            return -1;
        }
        if (x[0] > y[0]) {
            return 1;
        }
        return 0;
    });
    seriesVolume.setData(JSON.parse(JSON.stringify(volumeItems)));
    seriesQuantity.setData(JSON.parse(JSON.stringify(volumeQtde)));

}

function addVolumeQtde(  time , count )
{
    var found = false;

    for( var i=0;i<volumeQtde.length;i++ ){
        var item = volumeQtde[i];

        if( item[0]=== time ){
            if( item[1] == undefined ){
                item[1] = 0;
            }
            item[1] = item[1] + count;
            found = true;

            volumeQtde[i] = item;
            break;
        }
    }
    if( !found ){
        volumeQtde[volumeQtde.length] = [ time, count ];
    }
}

function addVolumeItem(  time, amount )
{
    console.log("$ " + time + " " + amount );

    var found = false;

    for( var i=0;i<volumeItems.length;i++ ){
        var item = volumeItems[i];

        if( item[0]== time ){
            if( item[1] == undefined ){
                item[1] = 0;
            }
            item[1] = item[1] + amount;
            found = true;
            break;
        }
    }
    if( !found ){
        volumeItems[volumeItems.length] = [ time, amount ];
    }
}



function createMonitorChart() {
    $(document).ready(function () {
        Highcharts.setOptions({
            global: {
                useUTC: false
            }
        });

        var chart  = $('#transactionVolume').highcharts('StockChart', {
            chart: {
                animation: Highcharts.svg, // don't animate in old IE
                marginRight: 10,
                events: {
                    load: function () {

                        seriesVolume = this.series[0];

                    }
                }
            },
            rangeSelector : {
                enabled: false
            },

            title: {
                text: 'Valor bruto'
            },
            xAxis: {
                type: 'datetime',
                tickPixelInterval: 150,
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
                enabled: false
            },
            exporting: {
                enabled: false
            },
            credits: {
                enabled: false

            },
            series: [{
                name: 'Transações',
                data: [],
                type: "line",
                dataGrouping:{
                    approximation: "sum",
                    enabled: false
                },
                fillColor : {
                    linearGradient : {
                        x1: 0,
                        y1: 0,
                        x2: 0,
                        y2: 1
                    },
                    stops : [
                        [0, Highcharts.getOptions().colors[0]],
                        [1, Highcharts.Color(Highcharts.getOptions().colors[0]).setOpacity(0).get('rgba')]
                    ]
                }
            }]
        });

        $('#transactionQuantity').highcharts('StockChart', {
            chart: {
                type: 'column',
                animation: Highcharts.svg, // don't animate in old IE
                marginRight: 10,
                events: {
                    load: function () {
                        seriesQuantity = this.series[0];
                    }
                }
            },
            rangeSelector : {
                enabled: false
            },
            navigator: {
                enabled: true
            },

            title: {
                text: 'Quantidade transações'
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
                enabled: false
            },
            exporting: {
                enabled: false
            },
            credits: {
                enabled: false

            },
            series: [{
                name: 'Transações',
                data: [],
                type: "line",
                pointInterval: 60 * 1000,
                dataGrouping:{
                    approximation: "sum",
                    units: [ [
                        'minute',
                        [1, 2, 5, 10, 15, 30]
                    ], [
                        'hour',
                        [1, 2, 3, 4, 6, 8, 12]
                    ], [
                        'day',
                        [1]
                    ], [
                        'week',
                        [1]
                    ], [
                        'month',
                        [1, 3, 6]
                    ], [
                        'year',
                        null
                    ]]
                },
                fillColor : {
                    linearGradient : {
                        x1: 0,
                        y1: 0,
                        x2: 0,
                        y2: 1
                    },
                    stops : [
                        [0, Highcharts.getOptions().colors[0]],
                        [1, Highcharts.Color(Highcharts.getOptions().colors[0]).setOpacity(0).get('rgba')]
                    ]
                }
            }]
        });
    });
}


