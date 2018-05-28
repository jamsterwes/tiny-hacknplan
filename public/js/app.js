var goWebApp = angular.module("goWebApp", [])

goWebApp.controller('HomeListController', function HomeListController($scope, $http) {
    $scope.tasks = [];
    $scope.username = "";

    $scope.refreshData = function() {
        if ($scope.username != "") {
            $http.get("/api/my_tasks/" + $scope.username)
            .then(function(response) {
                $scope.tasks = response.data;
            });
        } else {
            $scope.tasks = [];
        }
    }
})
