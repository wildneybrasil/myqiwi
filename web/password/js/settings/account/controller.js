techPayControllers.controller('MerchantEditCtrl',  function($scope, $location, $cookieStore, merchant) {
    var authToken = $cookieStore.get('authToken');




    var success = function (data) {
        switch (data.status) {
            case 'success':
                $scope.data = data.data;;
                console.log( JSON.stringify(data.data))
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



    //
    $scope.update = function(){
        var authToken = $cookieStore.get('authToken');

        var success = function (data) {
            switch (data.status) {
                case 'success':
                    alert( "Dados atualizados com sucesso" );
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
        merchant.update(authToken, $scope.data.name, $scope.data.email, $scope.data.cnpj, $scope.data.phone, $scope.data.address).success(success).error(error);
    }
});
