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

mgrApp.controller("saltkeymgrCtrl", function ($scope,$http,$modal,$log,
      $timeout,baseUrl,$rootScope) {

  $scope.environments = [];
  $scope.env = {};
  $scope.servernames = [];

  // Button enabling/disabling and content showing/hiding vars
  $scope.btnenvlistdisabled = false;
  $scope.btnlistkeysdisabled = true;
  $scope.btnapplydisabled = true;
  $scope.btnreviewdisabled = true;
  $scope.envchosen = false;
  $scope.listbtnpressed = false;
  $scope.serverlist_ready = false;
  $scope.serverlist_empty = true;

  // ----------------------------------------------------------------------
  $scope.Review = function() {
  // ----------------------------------------------------------------------
    $scope.btnlistkeysdisabled = true;
    $scope.btnenvlistdisabled = false;
    $scope.btnapplydisabled = false;
    $scope.btnreviewdisabled = true;
  };

  // ----------------------------------------------------------------------
  $scope.Apply = function() {
  // ----------------------------------------------------------------------
    $scope.btnlistkeysdisabled = true;
    $scope.btnenvlistdisabled = false;
    $scope.btnapplydisabled = true;
    $scope.btnreviewdisabled = false;
  };

  // ----------------------------------------------------------------------
  $scope.Restart = function() {
  // ----------------------------------------------------------------------
    $scope.btnlistkeysdisabled = true;
    $scope.btnenvlistdisabled = false;
    $scope.btnapplydisabled = true;
    $scope.btnreviewdisabled = true;
    $scope.listbtnpressed = false;
    $scope.serverlist_ready = false;
    $scope.serverlist_empty = false;
    $scope.envchosen = false;
  };

  // ----------------------------------------------------------------------
  $scope.ServerList = function() {
  // ----------------------------------------------------------------------
    $scope.btnlistkeysdisabled = true;
    $scope.btnenvlistdisabled = true;
    $scope.btnapplydisabled = true;
    $scope.btnreviewdisabled = true;
    $scope.listbtnpressed = true;
    $scope.serverlist_ready = false;
    $scope.serverlist_empty = false;
    $scope.envchosen = true;

    $scope.FillServerListTable();
  };

  // ----------------------------------------------------------------------
  $scope.envChoice = function( envobj, $event ) {
  // ----------------------------------------------------------------------
    $event.preventDefault();
    $event.stopPropagation();
    $scope.status.isopen = !$scope.status.isopen;
    $scope.envchosen = true;
    $scope.btnlistkeysdisabled = false;
    $scope.btnenvlistdisabled = true;
    $scope.env = envobj;
  };

  // ----------------------------------------------------------------------
  $scope.GetOutputLine = function( id ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/outputlines?job_id=" + id
    }).success( function(data, status, headers, config) {

      $scope.bob = $.parseJSON(data[0].Text);
      var servernames = $.parseJSON(data[0].Text);

      // Copy servernames array to objects
      for( i=0; i<servernames.length; ++i ) {
        key = Object.keys(servernames[i])
        $scope.servernames[i] = {
          Name: key[0],
          Selected: false
        };
      }

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

  // ----------------------------------------------------------------------
  $scope.PollForJobFinish = function( id,delay ) {
  // ----------------------------------------------------------------------
      $timeout( function() {
        $http({
          method: 'GET',
          url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
               + "/jobs?job_id=" + id
        }).success( function(data, status, headers, config) {
          job = data[0];
          if(job.Status == 0 || job.Status == 1 || job.Status == 4) {
            // Then retry every second
            $scope.PollForJobFinish(id,1000);
          }
          if(job.Status == 5) {
            $scope.serverlist_ready = true;
            $scope.serverlist_empty = false;
            $scope.btnlistkeysdisabled = false;

            $scope.GetOutputLine( id );
          }
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
      }, delay );
  };

  // ----------------------------------------------------------------------
  $scope.FillServerListTable = function() {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltnewserver/servers?env_id=" + $scope.env.Id
    }).success( function(data, status, headers, config) {
      $scope.PollForJobFinish(data.JobId,100);
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

  // ----------------------------------------------------------------------
  $scope.FillEnvironmentsTable = function() {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/envs?writeable=1"
    }).success( function(data, status, headers, config) {
      $scope.environments = data;
      if( data.length == 0 ) {
        $scope.serverlist_empty = true;
    }
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

  $scope.FillEnvironmentsTable();

});
