
techPayControllers.controller('LoginCtrl',  function($scope, $location, $cookieStore, authorization) {

    $scope.doLogin = function () {

        // var loading_screen = pleaseWait({
        //     backgroundColor: '#2277ee',
        //     loadingHtml: "<div class='sk-spinner sk-spinner-wave'><div class='sk-rect1'></div><div class='sk-rect2'></div><div class='sk-rect3'></div><div class='sk-rect4'></div><div class='sk-rect5'></div></div>"
        // });

        var credentials = {
            username: this.username,
            password: this.password
        };

        var success = function (data) {
//            loading_screen.finish();
//            alert( JSON.stringify( data ));

            switch( data.message ){
                case 'success':
                    $cookieStore.put('token', data.data.hash);

                    $location.path('/home');
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

        authorization.doLogin(credentials).success(success).error(error);
    };
});

