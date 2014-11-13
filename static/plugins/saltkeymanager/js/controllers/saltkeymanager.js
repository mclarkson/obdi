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
  $scope.keylist = {};
  $scope.keyfilter = "";
  $scope.env = {};
  $scope.status = {};  // For env chooser button

  // Alerting
  $scope.message = "";
  $scope.mainmessage = "";
  $scope.okmessage = "";
  $scope.login.error = false;

  // Button enabling/disabling and content showing/hiding vars
  $scope.envchosen = false;
  $scope.listbtnpressed = false;
  $scope.btnenvlistdisabled = false;
  $scope.showkeybtnblockhidden = false;
  $scope.btnshowkeysdisabled = true;
  $scope.keylist_ready = false;
  $scope.keylist_empty = true;
  $scope.keylist_accept_empty = true;
  $scope.keylist_reject_empty = true;
  $scope.keylist_unaccepted_empty = true;

  // ----------------------------------------------------------------------
  var clearMessages = function() {
  // ----------------------------------------------------------------------
    $scope.message = "";
    $scope.mainmessage = "";
    $scope.okmessage = "";
    $scope.login.error = false;
    $scope.error = false;
  }

  // ----------------------------------------------------------------------
  $scope.envChoice = function( envobj, $event ) {
  // ----------------------------------------------------------------------
    clearMessages();
    $event.preventDefault();
    $event.stopPropagation();
    $scope.status.isopen = !$scope.status.isopen;
    $scope.envchosen = true;
    $scope.btnshowkeysdisabled = false;
    $scope.btnenvlistdisabled = true;
    $scope.env = envobj;
  };

  // KEY MANAGEMENT

  // ----------------------------------------------------------------------
  $scope.Accept = function( name ) {
  // ----------------------------------------------------------------------
    $http({
      method: 'POST',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltkeymanager/saltkeys?hostname=" + name
           + "&type=accept"
           + "&env_id=" + $scope.env.Id,
    }).success( function(data, status, headers, config) {
      $scope.PollForJobFinish(data.JobId,100,0,$scope.GetKeyOutputLine);
    }).error( function(data,status) {
      if (status>=500) {
        $scope.login.errtext = "Server error.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status==401) {
        $scope.login.errtext = "Session expired.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status>=400) {
        clearMessages();
        $scope.message = "Server said: " + data['Error'];
        $scope.error = true;
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
  }

  // ----------------------------------------------------------------------
  $scope.Reject = function( name ) {
  // ----------------------------------------------------------------------
    $http({
      method: 'POST',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltkeymanager/saltkeys?hostname=" + name
           + "&type=reject"
           + "&env_id=" + $scope.env.Id,
    }).success( function(data, status, headers, config) {
      $scope.PollForJobFinish(data.JobId,100,0,$scope.GetKeyOutputLine);
    }).error( function(data,status) {
      if (status>=500) {
        $scope.login.errtext = "Server error.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status==401) {
        $scope.login.errtext = "Session expired.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status>=400) {
        clearMessages();
        $scope.message = "Server said: " + data['Error'];
        $scope.error = true;
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
  }

  // ----------------------------------------------------------------------
  $scope.Delete = function( name ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'DELETE',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltkeymanager/saltkeys/" + name
           + "?env_id=" + $scope.env.Id,
    }).success( function(data, status, headers, config) {
      $scope.PollForJobFinish(data.JobId,100,0,$scope.GetKeyOutputLine);
    }).error( function(data,status) {
      if (status>=500) {
        $scope.login.errtext = "Server error.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status==401) {
        $scope.login.errtext = "Session expired.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status>=400) {
        clearMessages();
        $scope.message = "Server said: " + data['Error'];
        $scope.error = true;
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
  }

  // ----------------------------------------------------------------------
  $scope.GetKeyOutputLine = function( id ) {
  // ----------------------------------------------------------------------
  // Actually don't bother with the result, just refresh the list
  // since we know Salt has finished now.

    $scope.FillKeyListTable();
  }

  // ----------------------------------------------------------------------
  $scope.GetKeyListOutputLine = function( id ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/outputlines?job_id=" + id
    }).success( function(data, status, headers, config) {

      $scope.showkeybtnblockhidden = true;

      var keylist = $.parseJSON(data[0].Text);

      $scope.keylist = keylist;

      if( $scope.keylist.minions.length == 0 ) {
        $scope.keylist_accept_empty = true;
      } else {
        $scope.keylist_accept_empty = false;
      }
      if( $scope.keylist.minions_pre.length == 0 ) {
        $scope.keylist_unaccepted_empty = true;
      } else {
        $scope.keylist_unaccepted_empty = false;
      }
      if( $scope.keylist.minions_rejected.length == 0 ) {
        $scope.keylist_reject_empty = true;
      } else {
        $scope.keylist_reject_empty = false;
      }

      $scope.keylist_ready = true;
      $scope.keylist_empty = false;

    }).error( function(data,status) {
      if (status>=500) {
        $scope.login.errtext = "Server error.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status==401) {
        $scope.login.errtext = "Session expired.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status>=400) {
        clearMessages();
        $scope.message = "Server said: " + data['Error'];
        $scope.error = true;
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

  // DC, ENV & VERSION

  // ----------------------------------------------------------------------
  $scope.GetServerSettingOutputLine = function( id ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/outputlines?job_id=" + id
    }).success( function(data, status, headers, config) {

      try {
        var result = $.parseJSON(data[0].Text);
      } catch (e) {
        clearMessages();
        $scope.message = "Error: " + e;
        $scope.message_jobid = id;
      }

      if( result.length == 0  ||
          typeof( result[$scope.changeversionview.saltid] ) == undefined ) {
        $scope.message = "The configuration did not complete.";
        $scope.message_jobid = id;
      }

      $scope.okmessage = "Server configuration was updated successfully.";

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
  $scope.ApplyGrain = function( saltid, grain, data ) {
  // ----------------------------------------------------------------------
  // Send { Grain:"version",Text:"0.1.2" }

    clearMessages();

    config = {};
    config.Grain = grain;
    config.Text = data;

    $http({
      method: 'POST',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltconfigserver/grains?salt_id=" + saltid
           + "&env_id=" + $scope.env.Id,
      data: config
    }).success( function(data, status, headers, config) {
      $scope.PollForJobFinish(
        data.JobId,
        50,
        0,
        $scope.GetServerSettingOutputLine);
    }).error( function(data,status) {
      if (status>=500) {
        $scope.login.errtext = "Server error.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status==401) {
        $scope.login.errtext = "Session expired.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status>=400) {
        clearMessages();
        $scope.message = "Server said: " + data['Error'];
        $scope.error = true;
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

  // ENV MANAGEMENT

  // ----------------------------------------------------------------------
  $scope.EnvConfig = function( name ) {
  // ----------------------------------------------------------------------
  }

  // ----------------------------------------------------------------------
  $scope.Restart = function() {
  // ----------------------------------------------------------------------
    clearMessages();
    $scope.envchosen = false;
    $scope.listbtnpressed = false;
    $scope.btnenvlistdisabled = false;
    $scope.showkeybtnblockhidden = false;
    $scope.btnshowkeysdisabled = true;
    $scope.keylist_ready = false;
    $scope.keylist_empty = true;
  };

  // ----------------------------------------------------------------------
  $scope.KeyList = function() {
  // ----------------------------------------------------------------------
    $scope.btnshowkeysdisabled = true;
    $scope.listbtnpressed = true;
    $scope.keylist_ready = false;
    $scope.keylist_empty = false;

    $scope.FillKeyListTable();
  };

  // ----------------------------------------------------------------------
  $scope.PollForJobFinish = function( id,delay,count,func ) {
  // ----------------------------------------------------------------------
      $timeout( function() {
        $http({
          method: 'GET',
          url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
               + "/jobs?job_id=" + id
        }).success( function(data, status, headers, config) {
          job = data[0];
          if(job.Status == 0 || job.Status == 1 || job.Status == 4) {
            if( count > 40 ) {
              clearMessages();
              $scope.message = "Job took too long. check job ID " +
                               + id + ", then try again.";
              $scope.message_jobid = job['Id'];
            } else {
              // Then retry: capped exponential backoff
              delay = delay < 600 ? delay * 2 : 1000;
              count = count + 1;
              $scope.PollForJobFinish(id,delay,count,func);
            }
          } else if(job.Status == 5) { // Job was successfully completed
            func( id );
          } else { // Some error
            clearMessages();
            $scope.message = "Server said: " + job['StatusReason'];
            $scope.message_jobid = job['Id'];
          }
        }).error( function(data,status) {
          if (status>=500) {
            $scope.login.errtext = "Server error.";
            $scope.login.error = true;
            $scope.login.pageurl = "login.html";
          } else if (status==401) {
            $scope.login.errtext = "Session expired.";
            $scope.login.error = true;
            $scope.login.pageurl = "login.html";
          } else if (status>=400) {
            clearMessages();
            $scope.message = "Server said: " + data['Error'];
            $scope.error = true;
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
  $scope.FillKeyListTable = function() {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltkeymanager/saltkeys"
           + "?env_id=" + $scope.env.Id
    }).success( function(data, status, headers, config) {
      $scope.PollForJobFinish(data.JobId,100,0,$scope.GetKeyListOutputLine);
    }).error( function(data,status) {
      if (status>=500) {
        $scope.login.errtext = "Server error.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status==401) {
        $scope.login.errtext = "Session expired.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status>=400) {
        clearMessages();
        $scope.message = "Server said: " + data['Error'];
        $scope.error = true;
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
        $scope.btnenvlistdisabled = true;
      }
    }).error( function(data,status) {
      if (status>=500) {
        $scope.login.errtext = "Server error.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status==401) {
        $scope.login.errtext = "Session expired.";
        $scope.login.error = true;
        $scope.login.pageurl = "login.html";
      } else if (status>=400) {
        clearMessages();
        $scope.mainmessage = "Server said: " + data['Error'];
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
