
techPayControllers.controller('LoginCtrl',  function($scope, $location, $cookieStore, $stateParams, authorization) {

    $scope.chpassword = function () {
	if( $scope.password1 != $scope.password2 ){
	    $scope.message = "As senhas diferem, verifique as senhas";
	    return ;
	}
        // var loading_screen = pleaseWait({
        //     backgroundColor: '#2277ee',
        //     loadingHtml: "<div class='sk-spinner sk-spinner-wave'><div class='sk-rect1'></div><div class='sk-rect2'></div><div class='sk-rect3'></div><div class='sk-rect4'></div><div class='sk-rect5'></div></div>"
        // });

        var credentials = {
            email: $stateParams.email,
            token: $stateParams.token,
	    password: $scope.password1
        };

        var success = function (data) {
//            loading_screen.finish();
//            alert( JSON.stringify( data ));

            switch( data.status ){
                case 'success':
		    alert("Sua senha foi alterada com sucesso");
                    break;
                default:
                    alert( JSON.stringify( data )  );
                    break;
            }
        };

        var error = function (v) {
            alert("ERRO " );
            // TODO: apply user notification here..
        };

        authorization.chpassword(credentials).success(success).error(error);
    };
    $scope.verifyToken = function () {
	
        // var loading_screen = pleaseWait({
        //     backgroundColor: '#2277ee',
        //     loadingHtml: "<div class='sk-spinner sk-spinner-wave'><div class='sk-rect1'></div><div class='sk-rect2'></div><div class='sk-rect3'></div><div class='sk-rect4'></div><div class='sk-rect5'></div></div>"
        // });

        var credentials = {
            email: $stateParams.email,
            token: $stateParams.token
        };

        var success = function (data) {
//            loading_screen.finish();
//            alert( JSON.stringify( data ));

            switch( data.status ){
                case 'success':
//                    $location.path('/passwordChangedOk');
//		    alert("OK");
                    break;
                default:
                    alert( JSON.stringify( data )  );
                    break;
            }
        };

        var error = function (v) {
            alert("ERRO " );
            // TODO: apply user notification here..
        };

        authorization.verifyToken(credentials).success(success).error(error);
    };

     $scope.verifyToken();
});

