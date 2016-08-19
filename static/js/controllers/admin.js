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

// Set a global so we know what interface we're in (admin or user)
window.interface = "admin";

mgrApp.controller("adminCtrl", function ($scope) {

  $scope.admin = {};
  $scope.admin.currentPage = "users.html";

  clearActive();
  $scope.admin.usersActive = "active";

  function clearActive() {
    $scope.admin.usersActive = "";
    $scope.admin.environmentsActive = "";
    $scope.admin.datacentresActive = "";
    $scope.admin.scriptsActive = "";
    $scope.admin.pluginsActive = "";
  }

  $scope.admin.showUsersPage = function() {
    clearActive();
    $scope.admin.usersActive = "active";
    $scope.admin.currentPage = "users.html";
  }

  $scope.admin.showEnvironmentsPage = function() {
    clearActive();
    $scope.admin.environmentsActive = "active";
    $scope.admin.currentPage = "environments.html";
  }

  $scope.admin.showDatacentresPage = function() {
    clearActive();
    $scope.admin.datacentresActive = "active";
    $scope.admin.currentPage = "datacentres.html";
  }

  $scope.admin.showScriptsPage = function() {
    clearActive();
    $scope.admin.scriptsActive = "active";
    $scope.admin.currentPage = "scripts.html";
  }

  $scope.admin.showPluginsPage = function() {
    clearActive();
    $scope.admin.pluginsActive = "active";
    $scope.admin.currentPage = "plugins.html";
  }

  $scope.admin.page = function() {
    return $scope.admin.currentPage;
  }
});
