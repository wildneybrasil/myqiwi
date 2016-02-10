'use strict';

/* Services */

var techPayServices = angular.module('techPayServices', ['ngResource']);
var host = "http://ec2-54-207-24-178.sa-east-1.compute.amazonaws.com/ws"
//var host = "http://api.irisescola.com.br/wslogin"



techPayServices.factory('authorization', function ($http) {
  return {
    chpassword: function (credentials) {
      var jsonRequest = {
        email: credentials.email,
        password: credentials.password,
        token: credentials.token
      };

      return $http.post(host + '/chpassword', jsonRequest);
    }
  };
});

