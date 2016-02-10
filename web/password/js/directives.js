'use strict';

angular.module('techPayDirectives', [])
.controller('techPayControllers', ['$scope', function($scope) {


}])
 .directive('format', function ($filter) {
        return {
            require: 'ngModel',
            link: function (scope, element, attrs, ctrl) {
                if (!ctrl) return;

                element.bind("keyup", function () {
                    var result = $filter('dinheiro')(ctrl.$viewValue);
                    ctrl.$setViewValue(result);
                    ctrl.$render();
                })

            }
        };
    });