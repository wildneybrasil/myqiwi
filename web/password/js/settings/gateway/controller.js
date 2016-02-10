techPayControllers.controller('MerchantGateway',  function($scope, $location, $cookieStore, merchant) {
    var authToken = $cookieStore.get('authToken');

    var success = function (data) {
        switch (data.status) {
            case 'success':
                $scope.items = data.data.processorInfo;
                console.log( JSON.stringify(data))
                break;
            default:
                if( data.errorCode == 5020 ){
                    $location.path('/start');
                }
                alert( data.errorMessage  );
                break;
        }
    };

    var error = function () {
        // TODO: apply user notification here..
    };
    merchant.listProcessor(authToken).success(success).error(error);


    $scope.change = function( id ){
        for( var i=0;i<$scope.items.length;i++){
            $scope.items[i].main=0;
            if( $scope.items[i].id == id ){
                $scope.items[i].main=1;
            }
        }
    };


    $scope.update = function(){
        var success = function (data) {
            switch (data.status) {
                case 'success':
                    $scope.items = data.data.processorInfo;
                    console.log( JSON.stringify(data))
                    break;
                default:
                    if( data.errorCode == 5020 ){
                        $location.path('/start');
                    }
                    alert( data.errorMessage  );
                    break;
            }
        };

        var error = function () {
            // TODO: apply user notification here..
        };

        for( var i=0;i<$scope.items.length;i++){
            if( $scope.items[i].main == '1' ){
                console.log("UPDATING " +  JSON.stringify($scope.items[i]) );
                merchant.updateProcessor(authToken, $scope.items[i].id, $scope.items[i].main, $scope.items[i]).success(success).error(error);

                break;
            }
        }
    }
});
