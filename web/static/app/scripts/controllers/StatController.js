'use strict';

angular.module('ircbotApp').controller('StatController', ['$scope', 'StatProvider', function ($scope, StatProvider) {
    var colors = ['#1BE7FF', '#6EEB83', '#E4FF1A', '#E8AA14', '#FF5714', '#50514F', '#F25F5C', '#247BA0', '#70C1B3'];
    var shuffle = function(o) {
        for(var j, x, i = o.length; i; j = Math.floor(Math.random() * i), x = o[--i], o[i] = o[j], o[j] = x);
        return o;
    }
    shuffle(colors);
    StatProvider.get().$promise.then(function(data) {
        var linkData = [];
        var userData = [];
        var duplicateData = [];
        for (var i = 0; i < data.GroupLink.length; i++) {
            var e = data.GroupLink[i]
            linkData.push({value: e.Count, label: e.Type.capitalizeFirstLetter(),
                color: getColor(i)})
        }
        for (var i = 0; i < data.GroupUser.length; i++) {
            var e = data.GroupUser[i];
            userData.push({value: e.Count, label: e.User_name.capitalizeFirstLetter(),
                color: getColor(i + data.GroupLink.length)})
        }
        for (var i = 0; i < data.Duplicates.length; i++) {
            var e = data.Duplicates[i];
            duplicateData.push({value: e.Count, label: e.User_name.capitalizeFirstLetter(),
                color: getColor(i + data.GroupLink.length + data.GroupUser.length)})
        }
        var ctx = document.getElementById("linkCanvas").getContext("2d");
        var ctx2 = document.getElementById("userCanvas").getContext("2d");
        var ctx3 = document.getElementById("duplicateCanvas").getContext("2d");
        var linkChart = new Chart(ctx).Doughnut(linkData);
        var userChart = new Chart(ctx2).Doughnut(userData);
        var duplicateChart = new Chart(ctx3).Doughnut(duplicateData);
    });

    var getColor = function(i) {
        if (colors.length <= i) {
            var start = i - colors.length * Math.floor(i / colors.length);
            return colors[start];
        }
        return colors[i];
    }
}]);
