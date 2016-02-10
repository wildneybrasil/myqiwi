techPayControllers.controller('MerchantPayment',  function($scope, $location, $cookieStore, merchant) {
    var authToken = $cookieStore.get('authToken');

    var success = function (data) {
        switch (data.status) {
            case 'success':
                $scope.urlCallback = data.data.callback_url;;
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
    merchant.read(authToken).success(success).error(error);




    $scope.update = function(){
        var authToken = $cookieStore.get('authToken');

        var success = function (data) {
            switch (data.status) {
                case 'success':
                    $scope.data = data.data;;
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
        merchant.callbackURL(authToken, $scope.urlCallback ).success(success).error(error);
    }
});

