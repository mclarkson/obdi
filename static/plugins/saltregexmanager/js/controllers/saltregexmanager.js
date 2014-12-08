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

mgrApp.controller("saltregexmgrCtrl", function ($scope,$http,$modal,$log,
      $timeout,baseUrl,$rootScope) {

  $scope.environments = [];
  $scope.regexlist = {};
  $scope.keyfilter = "";
  $scope.mapfilter = "";
  $scope.env = {};
  $scope.status = {};  // For env chooser button
  $scope.forminvalid = true; // For grains setting (dc,env,ver)

  // Alerting
  $scope.message = "";
  $scope.mainmessage = "";
  $scope.okmessage = "";
  $scope.login.error = false;

  // Button enabling/disabling and content showing/hiding vars
  $scope.envchosen = {};
  $scope.envchosen.shown = false;
  $scope.mapconfig = {};
  $scope.mapconfig.map = [];
  $scope.mapconfig.shown = false;
  $scope.mapconfig.maplist_empty = true;
  $scope.mapconfig.maplist_ready = false;
  $scope.mapconfig.saltid = "";
  $scope.mapconfig.regx_name = "";
  $scope.listbtnpressed = false;
  $scope.btnenvlistdisabled = false;
  $scope.showkeybtnblockhidden = false;
  $scope.btnshowkeysdisabled = true;
  $scope.regexlist_ready = false;
  $scope.regexlist_empty = true;

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
  $scope.copyToController = function( isit ) {
  // ----------------------------------------------------------------------
      $scope.forminvalid = isit;
  }

  // ----------------------------------------------------------------------
  $scope.envChoice = function( envobj, $event ) {
  // ----------------------------------------------------------------------
    clearMessages();
    $event.preventDefault();
    $event.stopPropagation();
    $scope.status.isopen = !$scope.status.isopen;
    $scope.envchosen.shown = true;
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

  // DC, ENV & VERSION

  // ----------------------------------------------------------------------
  $scope.GetServerSettingOutputLine = function( id ) {
  // ----------------------------------------------------------------------

    //$scope.okmessage = "Server configuration was updated successfully.";
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
          typeof( result[$scope.envsetting.saltid] ) == undefined ) {
        $scope.message = "The configuration did not complete.";
        $scope.message_jobid = id;
      }

      $scope.envsetting.numupdated += 1;
      if( $scope.envsetting.numupdated == 3 ) {
        $scope.okmessage = "Server configuration was updated successfully.";
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
  $scope.ApplyGrain = function( saltid, grain, data ) {
  // ----------------------------------------------------------------------
  // Send { Grain:"version",Text:"0.1.2" }

    var config = {};
    config.Grain = grain;
    config.Text = data;

    $http({
      method: 'POST',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltconfigserver/grains?salt_id=" + saltid
           + "&env_id=" + $scope.env.Id,
      data: config
    }).success( function(data, status, headers, config) {
      $scope.PollForJobFinish( data.JobId, 50, 0,
        $scope.GetServerSettingOutputLine );
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
  $scope.ApplySettings = function() {
  // ----------------------------------------------------------------------

    clearMessages();

    $scope.envsetting.numupdated = 0;

    if( $scope.envsetting.dc ) {
      $scope.ApplyGrain( $scope.envsetting.saltid, "dc",
          $scope.envsetting.dc );
    } else {
      $scope.envsetting.numupdated += 1;
    }

    if( $scope.envsetting.env ) {
      $scope.ApplyGrain( $scope.envsetting.saltid, "env",
          $scope.envsetting.env );
    } else {
      $scope.envsetting.numupdated += 1;
    }

    if( $scope.envsetting.version ) {
      $scope.ApplyGrain( $scope.envsetting.saltid, "version",
          $scope.envsetting.version );
    } else {
      $scope.envsetting.numupdated += 1;
    }
  }

  // ENV MANAGEMENT

  // ----------------------------------------------------------------------
  $scope.MapConfig = function( regex_id, name ) {
  // ----------------------------------------------------------------------

    clearMessages();
    $scope.envchosen.shown = false;
    $scope.mapconfig.shown = true;
    $scope.mapconfig.saltid = regex_id;
    $scope.mapconfig.regx_name = name;

    $scope.FillMapsTable( regex_id );
  }

  // ----------------------------------------------------------------------
  $scope.FillMapsTable = function( regex_id ) {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltregexmanager/regex_sls_maps?regex_id=" + regex_id
           + "&env_id=" + $scope.env.Id
    }).success( function(data, status, headers, config) {

      // Extract data into array
      try {
        $scope.mapconfig.map = $.parseJSON(data.JsonData);
      } catch (e) {
        clearMessages();
        $scope.message = "Error: " + e;
        $scope.message_jobid = id;
      }

      if( $scope.mapconfig.map.length == 0 ) {
        $scope.mapconfig.maplist_empty = true;
      } else {
        $scope.mapconfig.maplist_empty = false;
      }

      // Let the view know

      $scope.mapconfig.maplist_ready = true;

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
  $scope.GoBack = function( ) {
  // ----------------------------------------------------------------------
    clearMessages();
    $scope.envchosen.shown = true;

    $scope.mapconfig.map = [];
    $scope.mapconfig.shown = false;
    $scope.mapconfig.maplist_empty = true;
    $scope.mapconfig.maplist_ready = false;
    $scope.mapconfig.saltid = "";
    $scope.mapconfig.regx_name = "";

    $rootScope.$broadcast( "setsearchtext", $scope.hostfilter );
  }

  // ----------------------------------------------------------------------
  $scope.Restart = function() {
  // ----------------------------------------------------------------------
    clearMessages();
    $scope.envchosen.shown = false;
    $scope.listbtnpressed = false;
    $scope.btnenvlistdisabled = false;
    $scope.showkeybtnblockhidden = false;
    $scope.btnshowkeysdisabled = true;
    $scope.keylist_ready = false;
    $scope.keylist_empty = true;
  };

  // ----------------------------------------------------------------------
  $scope.RegexList = function() {
  // ----------------------------------------------------------------------
    $scope.btnshowkeysdisabled = true;
    $scope.listbtnpressed = true;
    $scope.keylist_ready = false;
    $scope.keylist_empty = false;

    $scope.FillRegexListTable();
  };

  // ----------------------------------------------------------------------
  $scope.showOutputlines = function( id ) {
  // ----------------------------------------------------------------------
  // Redirect the user to the Jobs->Outputlines plugin

    $rootScope.outputlines_plugin = {};
    $rootScope.outputlines_plugin.id = id;
    $scope.setView( "plugins/jobs/html/outputlines.html" );
  }

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
  $scope.FillRegexListTable = function() {
  // ----------------------------------------------------------------------

    $http({
      method: 'GET',
      url: baseUrl + "/" + $scope.login.userid + "/" + $scope.login.guid
           + "/saltregexmanager/regexes"
           + "?env_id=" + $scope.env.Id
    }).success( function(data, status, headers, config) {

      $scope.showkeybtnblockhidden = true;

      var regexlist = $.parseJSON(data.JsonData);

      $scope.regexlist = regexlist;

      if( $scope.regexlist.length == 0 ) {
        $scope.regexlist_empty = true;
      } else {
        $scope.regexlist_empty = false;
      }

      $scope.regexlist_ready = true;

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

  // --------------------------------------------------------------------
  $scope.dialog = function (servername) {
  // --------------------------------------------------------------------

    $scope.servername = servername;

    var modalInstance = $modal.open({
      templateUrl: 'myModalContent.html',
      controller: $scope.ModalInstanceCtrl,
      size: 'sm',
      resolve: {
        // the servername variable is passed to the ModalInstanceCtrl
        servername: function () {
          return $scope.servername;
        }
      }
    });

    modalInstance.result.then(function (servername) {
      $log.info('Will delete: ' + $scope.servername + '(' + $scope.servername + ')' );
      $scope.Delete($scope.servername);
    }, function () {
      $log.info('Modal dismissed at: ' + new Date());
    });
  };

  // --------------------------------------------------------------------
  $scope.ModalInstanceCtrl = function ($scope, $modalInstance, servername) {
  // --------------------------------------------------------------------

    // So the template can access 'servername' in this new scope
    $scope.servername = servername;

    $scope.ok = function () {
      $modalInstance.close();
    };

    $scope.cancel = function () {
      $modalInstance.dismiss('cancel');
    };
  };

});
