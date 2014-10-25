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

//var mgrApp = angular.module("mgrApp", []);

mgrApp.constant("baseUrl", "/api");

mgrApp.controller("loginCtrl", function ($scope, $http, baseUrl) {

  $scope.login = {};
  $scope.login.pageurl = "login.html";
  $scope.login.error = false;
  $scope.login.errtext = "";

  // ------------------------------------------------------------------------
  $scope.clearlogin = function() {
  // ------------------------------------------------------------------------
    $scope.login.userid = "";
    $scope.login.password = "";
    $scope.login.guid = "";
  }

  // ------------------------------------------------------------------------
  $scope.dologin = function() {
  // ------------------------------------------------------------------------
    // Fix for browser autocomplete not registering
    $('input').checkAndTriggerAutoFillEvent();

    if( $scope.login.userid != "admin" && window.interface == "admin" ) {
      $scope.login.userid = '';
      $scope.login.password = '';
      $scope.login.error = true;
      $scope.login.errtext = "You do not have admin access.";
      $scope.login.pageurl = "login.html";
      return
    }

    if( $scope.login.userid == "admin" && window.interface == "run" ) {
      $scope.login.userid = '';
      $scope.login.password = '';
      $scope.login.error = true;
      $scope.login.errtext = "Admin user is not allowed here.";
      $scope.login.pageurl = "login.html";
      return
    }

    creds = {
      Login: $scope.login.userid,
      Password: $scope.login.password
    };

    $http({
      url: baseUrl + "/" + 'login',
      method: "POST",
      data: creds
    }).success( function (data) {
      $scope.login.guid = data.GUID;
      creds = {};
      $scope.login.password = '';
      if ($scope.login.userid == "admin") {
        $scope.login.pageurl = "admin.html";
      } else {
        $scope.login.pageurl = "run.html";
      }
    }).error( function(data,status) {
      $scope.login.userid = '';
      $scope.login.password = '';
      $scope.login.error = true;
      if (status>=500) {
        $scope.login.errtext = "Server error.";
      } else if (status>=400) {
        $scope.login.errtext = "Invalid Login or Password.";
      } else if (status==0) {
        // This is a guess really
        $scope.login.errtext = "Could not connect to server.";
      }
    });
  }

  // ------------------------------------------------------------------------
  $scope.Logout = function() {
  // ------------------------------------------------------------------------

    $http({
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
         + "/logout",
      method: "POST",
      data: creds
    }).success( function (data) {
      $scope.login.guid = '';
      $scope.login.userid = '';
      $scope.login.password = '';
      $scope.login.error = false;
      $scope.login.errtext = "";
      $scope.login.pageurl = "login.html";
    }).error( function(data,status) {
      $scope.login.guid = '';
      $scope.login.userid = '';
      $scope.login.password = '';
      $scope.login.error = true;
      $scope.login.pageurl = "login.html";
      if (status>=500) {
        $scope.login.errtext = "Server error.";
      } else if (status==401) {
        $scope.login.errtext = "Session expired.";
      } else if (status>=400) {
        $scope.login.errtext = "Invalid Login or Password.";
      } else if (status==0) {
        // This is a guess really
        $scope.login.errtext = "Could not connect to server.";
      }
    });
  }

  // ------------------------------------------------------------------------
  $scope.userOrAdminPage = function() {
  // ------------------------------------------------------------------------
    return $scope.login.pageurl;
  }
});

