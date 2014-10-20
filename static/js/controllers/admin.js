// Deployment Manager - a REST interface and GUI for deploying software
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
