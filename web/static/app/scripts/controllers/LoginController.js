'use strict';

 angular.module('ircbotApp').controller('LoginController', ['$scope', '$location', '$http', '$window', function ($scope, $location, $http, $window) {
	$scope.login = function(password) {
		$scope.loginError = "";
		$scope.loading = true;
		$http({
            method: 'POST',
            url: loginUrl,
            data: {
                password: password
            }
        }).success(function (data, status, headers, config) {
        	$window.sessionStorage.token = data.token;
        	$location.path('/');
        }).error(function(error) {
        	$scope.loginError = error.Error;
			$scope.loading = false;
        });
	};
 }]);
