var goWebApp = angular.module("goWebApp", [])

goWebApp.controller('HomeListController', function HomeListController($scope, $http) {
    $scope.tasks = [];
    $scope.users = [];
    $scope.showName = false;
    $scope.showComplete = true;
    $scope.showIncomplete = true;
    $scope.username = "";

    // task filters
    $scope.setShowAll = function() {
        $scope.showComplete = true;
        $scope.showIncomplete = true;
    }

    $scope.setShowComplete = function() {
        $scope.showComplete = true;
        $scope.showIncomplete = false;
    }

    $scope.setShowIncomplete = function() {
        $scope.showComplete = false;
        $scope.showIncomplete = true;
    }

    $scope.checkShowTask = function(Stage) {
        return (Stage.Name == 'Completed' && $scope.showComplete) || (Stage.Name != 'Completed' && $scope.showIncomplete);
    }

    // user filters
    $scope.setShowName = function() {
        $scope.showName = true;
    }

    $scope.setShowUsername = function() {
        $scope.showName = false;
    }

    $scope.getUsers = function() {
        $http.get("/api/users")
        .then(function(response) {
            $scope.users = response.data;
        });
    }

    $scope.loadUser = function(username) {
        $scope.username = username;
        $scope.refreshData();
    }

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

    $scope.getUsers();
})
