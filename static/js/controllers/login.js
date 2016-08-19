// Copyright 2016 Mark Clarkson
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

