'use strict';

angular.module('goPertApp')
  .controller('MainCtrl', function ($scope) {
    $scope.awesomeThings = [
      'HTML5 Boilerplate',
      'AngularJS',
      'Karma'
    ];
    $scope.SingPert = { m: 2, n: 2, l: { x: 1e-6, y: 0}} 
  });
