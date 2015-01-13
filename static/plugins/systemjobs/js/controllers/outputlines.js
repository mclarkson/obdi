// Obdi - a REST interface and GUI for deploying software
// Copyright (C) 2014  Mark Clarkson
// 
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// ------------------------------------------------------------------------
// AngularJS Controller
// ------------------------------------------------------------------------

mgrApp.controller("outputlinesCtrl", function ($scope,$http,$modal,$log,
  $timeout,baseUrl,$rootScope) {

  $scope.outputlines = [];

  $rootScope.$broadcast( "searchdisabled", true );

  // ----------------------------------------------------------------------
  $scope.goBack = function( id ) {
  // ----------------------------------------------------------------------

    $scope.setView( "plugins/systemjobs/html/view.html" );
  }

  // ----------------------------------------------------------------------
  var clearMessages = function() {
  // ----------------------------------------------------------------------
    $scope.message = "";
    $scope.okmessage = "";
    $scope.login.error = false;
    $scope.error = false;
  }

  // ----------------------------------------------------------------------
  $scope.FillOutputlinesTable = function() {
  // ----------------------------------------------------------------------

    id = $rootScope.outputlines_plugin.id;

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/outputlines?job_id=" + id
    }).success( function(data, status, headers, config) {
      $scope.outputlines = data;
    }).error( function(data,status) {
      if (status>=500) {
        $scope.login.errtext = "Server error.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status>=400) {
        $scope.login.errtext = "Session expired.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status==0) {
        // This is a guess really
        $scope.login.errtext = "Could not connect to server.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else {
        $scope.login.errtext = "Logged out due to an unknown error.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      }
    });
  };

  $scope.FillOutputlinesTable( );

});

