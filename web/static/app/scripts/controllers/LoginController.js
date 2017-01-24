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
        }).then(function(response) {
        	$window.sessionStorage.token = response.data.token;
        	$location.path('/');
        }, function(error) {
        	$scope.loginError = error.Error;
			$scope.loading = false;
        });
	};
 }]);
