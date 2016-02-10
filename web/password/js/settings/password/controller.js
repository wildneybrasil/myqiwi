techPayControllers.controller('MerchantPassword',  function($scope, $location, $cookieStore, merchant) {
    $scope.updatePassword = function(){
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
        if( $scope.newPassword1!=$scope.newPassword2){
            alert("A confirmação de senha não é a mesma que a senha")
        } else {
            merchant.password(authToken, $scope.oldPassword, $scope.newPassword1).success(success).error(error);
        }
    }
});

