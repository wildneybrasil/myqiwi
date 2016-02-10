'use strict';

/* App Module */

var techPayApp = angular.module('techPayApp', [
    'ngRoute',
    'ngCookies',
    'ui.router',
    'techPayControllers',
    'techPayServices',
    'techPayFilters',
    'techPayDirectives',
    'ngAnimate',
    'ng-currency'
]);
var techPayControllers = angular.module('techPayControllers', []);


techPayApp.config(function ($stateProvider, $urlRouterProvider) {
    $urlRouterProvider.otherwise("/password");

    $stateProvider
        .state('password', {
            url: "/password",
            templateUrl: "/password/password/password.html",
            controller: 'LoginCtrl'
        })

});

