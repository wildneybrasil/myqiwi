techPayControllers.controller('MainCtrl',  function($scope, $location, $cookieStore, authorization) {
    if( !$cookieStore.get('authToken') ){
        $location.path('/');

        return ;
    }


    var name = $cookieStore.get('name');
    $scope.name = name;

    $scope.logout = function(){
        $cookieStore.remove('authToken');
        $cookieStore.remove('merchantId');
        $cookieStore.remove('name');

        $location.path('/');
    }
});
