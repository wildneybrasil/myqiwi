techPayControllers.controller('HistoryCtrl',  function($scope, $location, $cookieStore, report) {
    var divLoading = angular.element( document.querySelector( '#loadingContainer' ) );

    $scope.main = false;

    $scope.load = function() {
        var authToken = $cookieStore.get('authToken');

        var success = function (data) {
            divLoading.remove(); // remove a div de animacao do DOM
            $scope.main = true;

            switch (data.status) {
                case 'success':
                    $scope.items = data.data;
                    //console.log(JSON.stringify(data))
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
            alert("Application error")
            $location.path('/start');
        };
        report.list(authToken, $scope.purchaseCode, $scope.selectedBrand, 50).success(success).error(error);
    };
    $scope.load();
    setupFilters();
});